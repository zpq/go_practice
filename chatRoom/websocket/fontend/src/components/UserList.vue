<template lang="html">
    <div id="userList">
        <ul>
            <li v-for="use in userLists">
                {{ user.name }}
            </li>
        </ul>
    </div>
</template>

<script>
    export default {
      data: function () {
        return {
            room : {
                Id : 1
            },
            userLists : {}
        }
      },
      computed: {},
      ready: function () {
            let url = "http://localhost:8008/room/:id/users";
            // $.getJSON(url, {roomId : this.room.Id}, function(data, status, xhr) {
            //     this.$set("userlists", JSON.Parse(data))
            // });
      },
      attached: function () {},
      methods: {},
      components: {},
      props : ['roomId']
    }
</script>

<style lang="css" scoped>
    #userList {
        margin : 10px;
        margin-left: 40px;
        display: inline-block;
        min-width: 200px;
        min-height: 500px;
        border:1px solid black;
        border-radius: 5px;
    }

    #userList ul li {
        list-style-type: none;
        padding: 10px;
    }
</style>
