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
          var _vm = this;
          $.ajax({
            url : "http://localhost:8008/login",
            type : 'POST',
            data : {username : this.username, password : this.password},
            dataType : "json",
            success : function(res) {
                console.log(res)
                if (res.status == 1) {
                    window.localStorage.setItem('chatRoom_username', res.data[0].username)
                    window.localStorage.setItem('chatRoom_token', res.data[0].token)
                    _vm.$route.router.go({name:'list'})
                } else {
                    alert(res.message)
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
