package webauthn

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"

	harpocrates "github.com/TunnelWork/Harpocrates"
	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/duo-labs/webauthn/protocol"
	duo "github.com/duo-labs/webauthn/webauthn"
)

type WebAuthn struct {
	duoWebAuthn *duo.WebAuthn
	config      duo.Config
}

func NewWebAuthn(conf map[string]string) *WebAuthn {
	RDDisplayName, ok := conf["RDDisplayName"]
	if !ok {
		RDDisplayName = "Ulysses Unknown Displayname"
	}
	RPID, ok := conf["RPID"]
	if !ok {
		RPID = "localhost"
	}
	RPOriginURL, ok := conf["RPOriginURL"]
	if !ok {
		RPOriginURL = "https://" + RPID
	}
	RPIconURL, ok := conf["RPIconURL"]
	if !ok {
		RPIconURL = ""
	}
	config := duo.Config{
		RPDisplayName: RDDisplayName,
		RPID:          RPID,
		RPOrigin:      RPOriginURL,
		RPIcon:        RPIconURL,
	}
	duoWebAuthn, err := duo.New(&config)
	if err != nil {
		return nil
	}
	return &WebAuthn{
		duoWebAuthn: duoWebAuthn,
		config:      config,
	}
}

func (w *WebAuthn) Registered(userID uint64) bool {
	result, err := auth.MFAEnabled(userID, "webauthn")
	if err != nil {
		return false
	}
	return result
}

func (w *WebAuthn) InitSignUp(userID uint64, username string) (map[string]interface{}, error) {
	var user *User
	user, err := LoadUser(userID)
	if err != nil {
		user, err = CreateNewUser(userID, username, w.config.RPIcon)
		if err != nil {
			return nil, err
		}
	}

	options, sessionData, err := w.duoWebAuthn.BeginRegistration(
		user,
	)

	if err != nil {
		return nil, err
	}

	// save sessionData to tmp
	sessionDataJson, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}
	sessionKey, err := harpocrates.GetRandomHex(16)
	if err != nil {
		return nil, err
	}
	err = auth.InsertTmpEntry(userID, "webauthn", sessionKey, string(sessionDataJson))
	if err != nil {
		return nil, err
	}

	// make response
	return map[string]interface{}{
		"options":    options,    // Actual WebAuthn options
		"sessionKey": sessionKey, // Include this key for server to know which session to use
	}, nil
}

func (w *WebAuthn) CompleteSignUp(userID uint64, mfaConf map[string]string) error {
	user, err := LoadUser(userID)
	if err != nil {
		return err
	}

	// load session
	sessionKey, ok := mfaConf["sessionKey"]
	if !ok {
		return errors.New("webauthn: incomplete post form")
	}
	sessionDataJson, err := auth.ReadTmpEntry(userID, "webauthn", sessionKey)
	if err != nil {
		return err
	}
	sessionData := duo.SessionData{}
	err = json.Unmarshal([]byte(sessionDataJson), &sessionData)
	if err != nil {
		return err
	}

	// verify userID
	if string(sessionData.UserID) != string(user.WebAuthnID()) {
		return errors.New("webauthn: session data userID mismatch")
	}

	// load response
	response, ok := mfaConf["response"]
	if !ok {
		return errors.New("webauthn: incomplete post form")
	}

	// verify response
	var ccr protocol.CredentialCreationResponse
	err = json.Unmarshal([]byte(response), &ccr)
	if err != nil {
		return err
	}
	if ccr.ID == "" {
		return errors.New("webauthn: parse error for Registration - missing ID")
	}
	testB64, err := base64.RawURLEncoding.DecodeString(ccr.ID)
	if err != nil || !(len(testB64) > 0) {
		return errors.New("webauthn: parse error for Registration - ID not base64.RawURLEncoded")
	}
	if ccr.PublicKeyCredential.Credential.Type == "" {
		return errors.New("webauthn: parse error for Registration - missing Credential.Type")
	}
	if ccr.PublicKeyCredential.Credential.Type != "public-key" {
		return errors.New("webauthn: parse error for Registration - Credential.Type not public-key")
	}

	var pcc protocol.ParsedCredentialCreationData
	pcc.ID, pcc.RawID, pcc.Type = ccr.ID, ccr.RawID, ccr.Type
	pcc.Raw = ccr

	parsedAttestationResponse, err := ccr.AttestationResponse.Parse()
	if err != nil {
		return err
	}

	pcc.Response = *parsedAttestationResponse

	shouldVerifyUser := w.duoWebAuthn.Config.AuthenticatorSelection.UserVerification == protocol.VerificationRequired
	invalidErr := (&pcc).Verify(sessionData.Challenge, shouldVerifyUser, w.duoWebAuthn.Config.RPID, w.duoWebAuthn.Config.RPOrigin)
	if invalidErr != nil {
		return invalidErr
	}

	credential, err := duo.MakeNewCredential(&pcc)
	if err != nil {
		return err
	}

	user.AuthnCredentials = append(user.AuthnCredentials, *credential)
	err = user.UpdateDatabase()
	return err
}

func (w *WebAuthn) NewChallenge(userID uint64) (map[string]interface{}, error) {
	user, err := LoadUser(userID)
	if err != nil {
		return nil, err
	}

	options, sessionData, err := w.duoWebAuthn.BeginLogin(user)
	if err != nil {
		return nil, err
	}

	// save sessionData to tmp
	sessionDataJson, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}
	sessionKey, err := harpocrates.GetRandomHex(16)
	if err != nil {
		return nil, err
	}
	err = auth.InsertTmpEntry(userID, "webauthn", sessionKey, string(sessionDataJson))
	if err != nil {
		return nil, err
	}

	// make response
	return map[string]interface{}{
		"options":    options,    // Actual WebAuthn options
		"sessionKey": sessionKey, // Include this key for server to know which session to use
	}, nil
}

func (w *WebAuthn) SubmitChallenge(userID uint64, challengeResponse map[string]string) error {
	user, err := LoadUser(userID)
	if err != nil {
		return err
	}

	// load session
	sessionKey, ok := challengeResponse["sessionKey"]
	if !ok {
		return errors.New("webauthn: incomplete post form")
	}
	sessionDataJson, err := auth.ReadTmpEntry(userID, "webauthn", sessionKey)
	if err != nil {
		return err
	}
	sessionData := duo.SessionData{}
	err = json.Unmarshal([]byte(sessionDataJson), &sessionData)
	if err != nil {
		return err
	}

	if string(sessionData.UserID) != string(user.WebAuthnID()) {
		return errors.New("webauthn: session data userID mismatch")
	}

	// load response
	response, ok := challengeResponse["response"]
	if !ok {
		return errors.New("webauthn: incomplete post form")
	}

	// verify response
	var car protocol.CredentialAssertionResponse
	err = json.Unmarshal([]byte(response), &car)
	if err != nil {
		return err
	}
	if car.ID == "" {
		return errors.New("webauthn: parse error for Login - missing ID")
	}

	_, err = base64.RawURLEncoding.DecodeString(car.ID)
	if err != nil {
		return errors.New("webauthn: parse error for Login - ID not base64.RawURLEncoded")
	}
	if car.Type != "public-key" {
		return errors.New("webauthn: parse error for Login - Credential.Type not public-key")
	}

	var par protocol.ParsedCredentialAssertionData
	par.ID, par.RawID, par.Type = car.ID, car.RawID, car.Type
	par.Raw = car

	par.Response.Signature = car.AssertionResponse.Signature
	par.Response.UserHandle = car.AssertionResponse.UserHandle

	err = json.Unmarshal(car.AssertionResponse.ClientDataJSON, &par.Response.CollectedClientData)
	if err != nil {
		return err
	}

	err = par.Response.AuthenticatorData.Unmarshal(car.AssertionResponse.AuthenticatorData)
	if err != nil {
		return err
	}

	parsedResponse := &par

	// From: duo.WebAuthn.FinishLogin()

	// Step 1
	userCredentials := user.WebAuthnCredentials()
	var credentialFound bool
	if len(sessionData.AllowedCredentialIDs) > 0 {
		var credentialsOwned bool
		for _, userCredential := range userCredentials {
			for _, allowedCredentialID := range sessionData.AllowedCredentialIDs {
				if bytes.Equal(userCredential.ID, allowedCredentialID) {
					credentialsOwned = true
					break
				}
				credentialsOwned = false
			}
		}
		if !credentialsOwned {
			return errors.New("webauthn: user does not own any of the allowed credentials")
		}
		for _, allowedCredentialID := range sessionData.AllowedCredentialIDs {
			if bytes.Equal(parsedResponse.RawID, allowedCredentialID) {
				credentialFound = true
				break
			}
		}
		if !credentialFound {
			return errors.New("webauthn: credential ID not found")
		}
	}

	// Step 2
	userHandle := parsedResponse.Response.UserHandle
	if len(userHandle) > 0 {
		if !bytes.Equal(userHandle, user.WebAuthnID()) {
			return errors.New("webauthn: user handle does not match user ID")
		}
	}

	// Step 3
	var loginCredential duo.Credential
	var credidx int
	for idx, cred := range userCredentials {
		if bytes.Equal(cred.ID, parsedResponse.RawID) {
			loginCredential = cred
			credentialFound = true
			credidx = idx
			break
		}
		credentialFound = false
	}

	if !credentialFound {
		return errors.New("webauthn: credential ID not found")
	}

	shouldVerifyUser := sessionData.UserVerification == protocol.VerificationRequired

	rpID := w.duoWebAuthn.Config.RPID
	rpOrigin := w.duoWebAuthn.Config.RPOrigin

	// Handle steps 4 through 16
	validError := parsedResponse.Verify(sessionData.Challenge, rpID, rpOrigin, shouldVerifyUser, loginCredential.PublicKey)
	if validError != nil {
		return validError
	}

	// Handle step 17
	loginCredential.Authenticator.UpdateCounter(parsedResponse.Response.AuthenticatorData.Counter)

	// Save updated credential to database
	user.AuthnCredentials[credidx] = loginCredential
	err = user.UpdateDatabase()
	if err != nil {
		return err
	}

	return nil
}

func (w *WebAuthn) Remove(userID uint64) error {
	return auth.ClearMFA(userID, "webauthn")
}
