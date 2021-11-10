function InitAddUser() {
  username = $("#regUserName").val()
  userid = parseInt($("#regUserID").val())

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
    .then((credential) => {
      document.getElementById("Secret").innerText = credential.secret;
      new QRCode(document.getElementById("QRcode"), credential.url);
    })
    .done(function (success) {
      alert("Generated QR Code and Secret for "+ username + "!")
      return
    })
    .catch((error) => {
      console.log(error)
      alert("failed to register " + username)
    }
  )
}

function VerifyAddUser() {
  userid = parseInt($("#regUserID").val());
  code = $("#regCode").val();
  secret = document.getElementById("Secret").innerText;

  $.post(
    '/register/finish',
    JSON.stringify({
      "userID": userid,
      "code": code,
      "secret": secret
    }),
    function (data) {
      return data
    },
    'json')
    .done(function (success) {
      alert("Registered!")
      return
    })
    .catch((error) => {
      console.log(error)
      alert("failed to register.")
    }
  )
}

function Login() {
  userid = parseInt($("#loginUserID").val());
  code = $("#loginCode").val();

  $.post(
    '/login/finish',
    JSON.stringify({
      "userID": userid,
      "code": code
    }),
    function (data) {
      return data
    },
    'json')
    .done(function (success) {
      alert("Logged In!")
      return
    })
    .catch((error) => {
      console.log(error)
      alert("Failed to LogIn.")
    }
  )
}