import Vue from 'vue'
// import App from './App'
import VueRouter from 'vue-router'
import Vuex from 'vuex'
import Hello from './components/Hello'

Vue.use(VueRouter)


var router = new VueRouter()
router.map({
    '/index' : {
        name : 'index',
        component : Index
        subRoutes: {
          '/selectCard': {
              component: Bar
          },
          '/register' : {
              name : 'register',
              component : Register
          },
          '/login' : {
              name : 'login',
              component : Login
          },
          '/profile' : {
            name :'profile',
            component : Profile,
          },
          '/fight' : {

          }
    },
});

router.redirect({
    "*" : "/index"
});

var App = Vue.extend({})
router.start(App, '#app')
