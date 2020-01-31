<template>
  <v-container class="px-0">
    <v-layout>
      <v-avatar>
        <img
          :src="message.speaker_avatar"
          :alt="message.speaker_name"
        >
      </v-avatar>
      <v-flex class="bubble ml-4 py-2 px-2">
        <div class="caret" />
        <template v-if="message.category === 'PLAIN_TEXT' ">
          <div class="px-1">
            {{ message.text }}
          </div>
        </template>
        <template v-else-if="message.category === 'PLAIN_IMAGE' ">
          <v-img
            :lazy-src="'data:image/jpeg;base64,' + message.attachment.thumbnail"
            :src="message.attachment.view_url"
            aspect-ratio="1"
            max-width="500"
            max-height="200"
          />
        </template>
        <template v-else-if="message.category === 'PLAIN_VIDEO' ">
          <video
            class="video"
            controls
          >
            <source
              :src="message.attachment.view_url"
              :type="message.attachment.mime_type"
            >
          </video>
        </template>
        <template v-else-if="message.category === 'PLAIN_AUDIO' ">
          <audio
            class="audio"
            controls
          >
            <source
              :src="message.attachment.view_url"
              :type="message.attachment.mime_type"
            >
          </audio>
        </template>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>

export default {
  name: 'Message',
  props: {
    message: {
      type: Object
    }
  },
  data: () => ({

  }),


};
</script>

<style lang="scss" scoped>
  .bubble {
    background: white;
    position: relative;
    box-shadow: 0 0 10px rgba(0,0,0,0.04);
    border-radius: 7px;
    font-size: 14px;
    .caret {
      position: absolute;
      left: -20px;
      width: 0;
      height: 0;
      border-width: 8px 12px 8px 12px;
      border-color: transparent white transparent transparent;
      border-style: solid;
    }
    .video {
      width: 100%;
      max-height: 200px;
    }
  }
</style>