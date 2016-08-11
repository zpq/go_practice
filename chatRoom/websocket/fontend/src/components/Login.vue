<template lang="html">
    <div>
        username : <input type="text" v-model="username" > <br>
        password : <input type="password" v-model="password" > <br>
        <button @click="login">login</button>
    </div>
</template>

<script>
export default {
  data: function () {
    return {
        username : '',
        password : ''
    }
  },
  computed: {},
  ready: function () {},
  attached: function () {},
  methods: {
      login : function() {
          $.ajax({
            url : "http://localhost:8008/login",
            type : 'POST',
            data : {username : this.username, password : this.password},
            dataType : "json",
            success : function(res) {
                console.log(res)
                if (res.Status == 1) {
                    window.localStorage.setItem('chatRoom_username', res.data.Username)
                    window.localStorage.setItem('chatRoom_token', res.data.Token)
                    alert(res.data[0].Username)
                } else {
                    alert(res.Message)
                }
            },
            error : function (res) {
                console.log(res)
            }
          });


      }
  },
  components: {}
}
</script>

<style lang="css">
</style>
