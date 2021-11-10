// Base64 to ArrayBuffer
function bufferDecode(value) {
    return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}

  // ArrayBuffer to URLBase64
function bufferEncode(value) {
return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");;
}

function registerUser() {

sessionK = ""  
username = $("#userName").val()
userid = parseInt($("#userID").val())

if (username === "") {
    alert("Please enter a username");
    return;
}

$.post(
    '/register/init',
    JSON.stringify({
        "userID": userid,
        "userName": username
    }),
    function (data) {
        return data
    },
    'json')
    .then((credentialCreationOptions) => {
      options = credentialCreationOptions.options;
      sessionK = credentialCreationOptions.sessionKey;
      credentialCreationOptions = options;
      // console.log("credentialCreationOptions: ",credentialCreationOptions);
      credentialCreationOptions.publicKey.challenge = bufferDecode(credentialCreationOptions.publicKey.challenge);
      credentialCreationOptions.publicKey.user.id = bufferDecode(credentialCreationOptions.publicKey.user.id);
      if (credentialCreationOptions.publicKey.excludeCredentials) {
        for (var i = 0; i < credentialCreationOptions.publicKey.excludeCredentials.length; i++) {
          credentialCreationOptions.publicKey.excludeCredentials[i].id = bufferDecode(credentialCreationOptions.publicKey.excludeCredentials[i].id);
        }
      }

      return navigator.credentials.create({
          publicKey: credentialCreationOptions.publicKey
        });
    })
    .then((credential) => {
      // sessionK = credential.sessionK;
      // credential = credential.credCont;
      console.log("credential:", credential);
      let attestationObject = credential.response.attestationObject;
      let clientDataJSON = credential.response.clientDataJSON;
      let rawId = credential.rawId;

      $.post(
        '/register/finish',
        JSON.stringify(
        { 
          "userID": userid,
          'sessionKey': sessionK,
          'response': {
            id: credential.id,
            rawId: bufferEncode(rawId),
            type: credential.type,
            response: {
              attestationObject: bufferEncode(attestationObject),
              clientDataJSON: bufferEncode(clientDataJSON),
            },
          }
        }
        ),
        function (data) {
          return data
        },
        'json')
        .fail(function (error) {
          throw(error)
        })
    })
    .done(function (success) {
      alert("successfully registered " + username + "!")
      return
    })
    .catch((error) => {
      console.log(error)
      alert("failed to register " + username)
    })
}

function loginUser() {
sessionK = ""  
username = $("#userName").val()
userid = parseInt($("#userID").val())
if (username === "") {
    alert("Please enter a username");
    return;
}

$.post(
    '/login/init',
    JSON.stringify({
      "userID": userid
    }),
    function (data) {
    return data
    },
    'json')
    .then((credentialRequestOptions) => {
    sessionK = credentialRequestOptions.sessionKey;
    credentialRequestOptions = credentialRequestOptions.options;
    console.log(credentialRequestOptions)
    credentialRequestOptions.publicKey.challenge = bufferDecode(credentialRequestOptions.publicKey.challenge);
    credentialRequestOptions.publicKey.allowCredentials.forEach(function (listItem) {
        listItem.id = bufferDecode(listItem.id)
    });

    return navigator.credentials.get({
        publicKey: credentialRequestOptions.publicKey
    })
    })
    .then((assertion) => {
    console.log(assertion)
    let authData = assertion.response.authenticatorData;
    let clientDataJSON = assertion.response.clientDataJSON;
    let rawId = assertion.rawId;
    let sig = assertion.response.signature;
    let userHandle = assertion.response.userHandle;

    $.post(
        '/login/finish',
        JSON.stringify({
          "userID": userid,
          'sessionKey': sessionK,
          "response": {
            id: assertion.id,
            rawId: bufferEncode(rawId),
            type: assertion.type,
            response: {
                authenticatorData: bufferEncode(authData),
                clientDataJSON: bufferEncode(clientDataJSON),
                signature: bufferEncode(sig),
                userHandle: bufferEncode(userHandle),
            },
          },
        }),
        function (data) {
        return data
        },
        'json')
    })
    .then((success) => {
    alert("successfully logged in " + username + "!")
    return
    })
    .catch((error) => {
    console.log(error)
    alert("failed to register " + username)
    })
}