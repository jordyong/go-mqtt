
const clientID = 'mqttjs_' + Math.random().toString(16).substr(2, 8)
const host = 'ws://172.27.2.35:9001/mqtt'
const options = {
  keepalive: 30,
  clientID,
  protocolId: 'MQTT',
  protocolVersion: 4,
  clean: true,
  reconnectPeriod: 1000,
  connectTimeout: 30 * 1000,
  will: {
    topic: 'WillMsg',
    payload: 'Connection Closed abnormally..!',
    qos: 0,
    retain: false,
    manualConnect: true,
  },
  rejectUnauthorized: false
}

//const client = mqtt.connect(host, options)

client.on('error', (err) => {
  console.log(err)
  client.end()
})

client.on('connect', () => {
  console.log('client connected:' + clientID)
})

client.on('message', (topic, message, packet) => {
  console.log('Received Message:= ' + message.toString() + '\nOn topic:= ' + topic)
})

document.getElementById('connect').addEventListener('click', function() {
  if (this.textContent === "Connect") {
    button = this;
    client.connect()
    client.on('connect', function() {
      button.textContent = "Disconnect"
    })
  }
  else if (this.textContent === "Disconnect") {
    client.end()
    client.on('disconnect', function() {
      this.textContent = "Connect";
    })
  }
});
