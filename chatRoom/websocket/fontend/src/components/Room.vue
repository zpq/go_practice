<template lang="html">
    <div>
        <chat-panel room-id="this.room.id"></chat-panel>
        <user-list room-id="this.room.id"></user-list>
    </div>
</template>

<script>

    import ChatPanel from './ChatPanel'
    import UserList from './UserList'

    export default {
      data: function () {
        return {
            room : {
                Id : 1
            },
            userName : ''
        }
      },
      computed: {},
      ready: function () {
          // check auth
          let username = window.localStorage.getItem('chatRoom_username');
          if (username) {
              let token = window.localStorage.getItem('chatRoom_token');
            //   $.ajax({
            //       url : ""
            //   })
          } else {
              this.$route.router.go({name : "login"});
          }
      },
      attached: function () {},
      methods: {},
      components: {
          ChatPanel,
          UserList
      }
    }
</script>

<style lang="css" scoped>

</style>
