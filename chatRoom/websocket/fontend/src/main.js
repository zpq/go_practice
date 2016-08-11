import Vue from 'vue'
import App from './App'
import VueRouter from 'vue-router'
Vue.use(VueRouter)
var router = new VueRouter();

import RoomList from './components/RoomList'
import Room from './components/Room'
import Hello from './components/Hello'
import Login from './components/Login'

//define all routes of your app
router.map({
	'/index' : {
		name : 'index',
		component : Hello,
	},
	'/list' : {
		name : 'list',
		component : RoomList
	},
	'/room/:id' : {
		name : 'room',
		component : Room
	},
	'/login' : {
		name : 'login',
		component : Login
	}
});

//default route
router.redirect({
	"*" : "/index"
})

router.start(App, '#app')
