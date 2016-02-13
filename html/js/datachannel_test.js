
// Initialise DataChannel.js
var datachannel = new DataChannel();


datachannel.userid = userid;
var ws;
ws = new WebSocket(wspath);
ws.onopen = function(evt) {
  console.log("OPEN");
}
ws.onclose = function(evt) {
  console.log("CLOSE");
  ws = null;
}

ws.onerror = function(evt) {
  console.log("ERROR: " + evt.data);
}

ws.onmessage = function(evt) {
    var jsonObject = JSON.parse(evt.data);
    console.log(evt.data);
    //config.onmessage(jsonObject);
  }

datachannel.openSignalingChannel = function(config) {
  var channel = config.channel || this.channel || "default-channel";
  var xhrErrorCount = 0;

  var socket = {
    send: function(message) {
      ws.send(JSON.stringify(message));
    },
    channel: channel
  };
  ws.onmessage = function(evt) {
    var jsonObject = JSON.parse(evt.data);
    console.log(evt.data);
    config.onmessage(jsonObject);
    //Now join this one
  }
  if (config.onopen) {
    setTimeout(config.onopen, 1);
  }
  return socket;
}


var onCreateChannel = function() {
  var channelName = cleanChannelName(channelInput.value);

  if (!channelName) {
    console.log("No channel name given");
    return;
  }

  disableConnectInput();

  datachannel.open(channelName);
};

var onJoinChannel = function() {
  var channelName = cleanChannelName(channelInput.value);

  if (!channelName) {
    console.log("No channel name given");
    return;
  }

  disableConnectInput();

  // Search for existing data channels
  datachannel.connect(channelName);
};



