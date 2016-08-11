<template>
	<div>
		<div class="roomList">
			<div class="roomThumb" v-link="{name : 'room'}">room1</div>
			<div class="roomThumb">room2</div>
			<div class="roomThumb">room3</div>
			<div class="roomThumb">room4</div>
		</div>
	</div>
</template>

<script>

export default {
	data () {
		return {
			lists : [

			]
		}
	},
	ready() {
		//get room list from server
	},
	methods : {
		goIndex () {
			this.$route.router.go({name : "index"});
		}
	}
}

</script>

<style scoped>
	.roomThumb {
		width: 200px;
		height : 200px;
		display: inline-block;
		box-shadow: 10px #123456;
		border: 1px solid grey;
		text-align: center;
		cursor: pointer;
		vertical-align:middle;
	}
</style>
