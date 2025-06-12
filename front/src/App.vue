<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <RtcPlayer src="http://localhost:8083/play/hls/JSESSIONID=86F633899E8C2CF1577717/index.m3u8" width="640px" height="360px" />
  </div>
</template>

<script>
import RtcPlayer from './components/RtcPlayer';

export default {
  name: 'App',
  components: {
    RtcPlayer
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
