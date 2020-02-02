<template>
  <v-container>
    <template v-if="messages.length !== 0">
      <message-bubble
        v-for="(msg, index) in messages"
        :key="msg.id"
        :message="msg"
        :prev="getPrev(index, messages)"
        @click.native="handleClickMessage(msg)"
      />
    </template>
    <template v-else>
      <p
        class="text-center font-weight-light mt-4"
        style="width: 100%;"
      >
        No messages
      </p>
    </template>
    <message-viewer ref="viewer" />
  </v-container>
</template>

<script>
import { mapState } from 'vuex'
import MessageViewer from './message/MessageViewer'
import MessageBubble from './message/MessageBubble'

export default {
  name: 'ChatArea',
  components: {
    MessageBubble,
    MessageViewer
  },
  computed: {
    ...mapState('message', {
      messages: state => state.messages
    })
  },
  methods: {
    getPrev(index, messages) {
      return index === 0 ? null : messages[index-1]
    },
    handleClickMessage(message) {
      const viewer = this.$refs.viewer
      if (message.category === 'PLAIN_IMAGE') {
        viewer.show(message)
      }
    }
  }
};
</script>

