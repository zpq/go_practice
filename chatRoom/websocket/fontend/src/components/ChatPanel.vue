<template lang="html">
    <div id="chatPanel">
        <div id="displayMessage"></div>
        <div id="sendTools">
            <input type="text" v-model="message">
            <button type="button" @click="sendMessage">send</button>
        </div>
    </div>
</template>

<script>
export default {
      data: function () {
        return {
            room : {
                Id : 1
            },
            message : '',
            ws : null
        }
      },
      props : ['roomId'],
      computed: {},
      ready: function () {

          if (window.WebSocket) {

              let displayMessage = $("#displayMessage");

              var ws = new WebSocket("ws:\/\/localhost:8008\/echo")
              this.ws = ws;
              ws.onopen = function(ev) {
                  ws.send("hello!");
                  console.log("connect success!");
              }

              ws.close = function(ev) {
                  ws = null;
                  console.log("connect broken!");
              }

              ws.onerror = function(evt) {
                  console.log("ERROR: " + evt.data);
              }

              ws.onmessage = function(ev) {
                  displayMessage.append("<p>" + ev.data + "</p>");
                  displayMessage.scrollTop(displayMessage.height()); //使滚动条在最底部
              }

          } else {
              document.write("Sorry! Your browsers does not support websocket!");
          }

      },
      attached: function () {},
      methods: {
        sendMessage : function() {
            this.ws.send(this.message);
            this.message = '';
        }
      },
      components: {}
}
</script>

<style lang="css">
#chatPanel{
    display: inline-block;
}
#chatPanel div {
    width: 600px;
}
#displayMessage {
    height : 200px;
    border : 1px solid #abcdef;
    overflow-y:auto;
}
#sendTools {
    margin-top: 20px;
    height : 50px;
    padding-top: 20px;
}
#sendTools input {
    display: inline-block;
    width: 70%;
    margin-right: 15px;
}
#sendTools button {
    display: inline-block;
    width: 20%;
}

</style>
