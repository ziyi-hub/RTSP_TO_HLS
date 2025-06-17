<template>
  <div id="app">
    <button
        v-for="camera in cameras"
        :key="camera.code"
        @click="switchCamera(camera.code)"
        :class="{ active: currentCamera === camera.code }"
    >
      {{ camera.name }}
    </button>

    <div id="playerContainer" style="margin-top: 20px;">
      <video
          ref="videoPlayer"
          width="640"
          height="360"
          controls
          autoplay
          muted
      ></video>
      <p>{{ status }}</p>
    </div>

<!--    <RtcPlayer src="http://localhost:8083/play/hls/JSESSIONID=86F633899E8C2CF1577717/index.m3u8" width="640px" height="360px" />-->
  </div>
</template>

<script>
// import RtcPlayer from './components/RtcPlayer';
import Hls from 'hls.js';

export default {
  name: 'App',
  components: {
    // RtcPlayer
  },

  data() {
    return {
      cameras: [
        { code: "41010400001310001212#16#4c6a6bade5c74d66b71971f6f2670e61", name: "摄像头 1" },
        { code: "41010400001310001290#16#4c6a6bade5c74d66b71971f6f2670e61", name: "摄像头 2" },
      ],
      currentCamera: null,
      status: "",
      hls: null,
    };
  },
  methods: {
    switchCamera(cameraCode) {
      this.status = `正在加载摄像头 ${cameraCode} 视频流...`;
      this.currentCamera = cameraCode;

      fetch(`/play_by_camera/${cameraCode}`)
          .then((res) => {
            if (!res.ok) throw new Error("摄像头不存在");
            return res.json();
          })
          .then((data) => {
            const video = this.$refs.videoPlayer;

            if (this.hls) {
              this.hls.destroy();
              this.hls = null;
            }

            if (Hls.isSupported()) {
              this.hls = new Hls();
              this.hls.loadSource(data.hlsUrl);
              this.hls.attachMedia(video);
              this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
                video.play();
                this.status = `正在播放摄像头 ${cameraCode}`;
              });
            } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
              video.src = data.hlsUrl;
              video.addEventListener("loadedmetadata", () => {
                video.play();
                this.status = `正在播放摄像头 ${cameraCode}`;
              });
            } else {
              this.status = "您的浏览器不支持播放HLS流";
            }
          })
          .catch((err) => {
            this.status = "加载失败：" + err.message;
            console.error(err);
          });
    },
  },
  mounted() {
    // 默认加载第一个摄像头
    if (this.cameras.length > 0) {
      this.switchCamera(this.cameras[0].code);
    }
  },
};
</script>

<style scoped>
button {
  margin: 5px;
  padding: 10px;
  cursor: pointer;
}
button.active {
  background-color: #409eff;
  color: white;
}
</style>
