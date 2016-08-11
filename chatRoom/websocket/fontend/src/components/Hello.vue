<template>
  <div class="hello">
    <h3>{{ msg }}</h3>
  </div>
</template>

<script>

export default {
  data () {
    return {
      msg: 'Welcome!'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1,h2,h3 {
  color: #42b983;
}
</style>
