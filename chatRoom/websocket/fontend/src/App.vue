<template>
    <div>
        <div id="header">
            <div>
                <a v-link="{ name: 'list' }">view list</a> |
                <a v-link="{ name: 'index' }">go index</a>
            </div>
        </div>
        <div id="container">
            <router-view></router-view>
        </div>
    </div>
</template>

<script>

  export default {
      components:{
          //若还有更多的组件，只需要在import的情况下，以逗号分割，继续注册就好
      },
      methods :{
          goIndex () {//方法，定义路由跳转，注意这里必须使用this，不然报错
              this.$route.router.go({name:"index"});
          }
      },
      ready() {
        //   $.get("http://localhost:8008/getToken", function(res) {
        //       console.log(res)
        //   }, function(res) {
        //       console.log(res)
        //   });
      }
  }
</script>


<style>

html {
  height: 100%;
}

body {
  /*display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;*/
}

a {
    text-decoration: none;
    background : #eed;
    border: 1px solid pink;
}

#container {
    max-width: 880px;
    height : 100%;
    margin: 0 auto;
    padding: 10px;
}

#header {
    width: 100%;
    border: 2px solid #abcdef;
    box-shadow: 5px 5px 5px #888888;
    margin: 10px 0;
    padding: 10px;
}


/*#app {
  color: #2c3e50;
  margin-top: -100px;
  font-family: Source Sans Pro, Helvetica, sans-serif;
  text-align: center;
}

#app a {
  color: #42b983;
  text-decoration: none;
}*/

/*.logo {
  width: 100px;
  height: 100px
}*/
</style>
