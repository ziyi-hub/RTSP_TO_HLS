<template>
  <div class="video-container">
    <video
        ref="video"
        controls
        autoplay
        muted
        class="video"
        :style="{ width: width, height: height }"
    ></video>
  </div>
</template>

<script>
import Hls from 'hls.js';

export default {
  name: 'RtcPlayer',
  props: {
    src: {
      type: String,
      required: true, // 例如：http://localhost:8083/play/hls/H264_AAC/index.m3u8
    },
    width: {
      type: String,
      default: '100%',
    },
    height: {
      type: String,
      default: 'auto',
    },
  },
  mounted() {
    const video = this.$refs.video;

    if (Hls.isSupported()) {
      this.hls = new Hls();
      this.hls.loadSource(this.src);
      this.hls.attachMedia(video);
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      // Safari 原生支持
      video.src = this.src;
    } else {
      console.error('HLS is not supported in this browser');
    }
  },
  beforeDestroy() {
    if (this.hls) {
      this.hls.destroy();
    }
  },
};
</script>

<style scoped>
.video-container {
  display: flex;
  justify-content: center;
  align-items: center;
}
.video {
  background-color: black;
}
</style>
