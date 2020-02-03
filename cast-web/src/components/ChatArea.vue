<template>
  <v-container v-mutate.child="handleMutate">
    <template v-if="messages.length !== 0">
      <message-bubble
        v-for="(msg, index) in messages"
        :key="msg.id"
        :idx="index"
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
    <div
      ref="bottom"
      v-intersect="handleBottomIntersect"
    />
  </v-container>
</template>

<script>
import { mapState } from 'vuex'
import MessageViewer from './message/MessageViewer'
import MessageBubble from './message/MessageBubble'
import { mapMutations } from 'vuex'

export default {
  name: 'ChatArea',
  components: {
    MessageBubble,
    MessageViewer
  },
  data() {
    return {
      isBottomIntersect: false
    }
  },
  computed: {
    ...mapState('message', {
      messages: state => state.messages
    })
  },
  mounted() {
    this.scrollToBottom()
    this.$root.$on('CHECK_NEW_MESSAGE', () => {
      this.scrollToBottom()
    })
  },
  methods: {
    ...mapMutations('message', {
      setHasNewMessage: 'SET_HAS_NEW_MESSAGE'
    }),
    getPrev(index, messages) {
      return index === 0 ? null : messages[index-1]
    },
    handleClickMessage(message) {
      const viewer = this.$refs.viewer
      if (message.category === 'PLAIN_IMAGE') {
        viewer.show(message)
      }
    },
    handleBottomIntersect(entries, observer) {
      this.isBottomIntersect = entries[0].isIntersecting
      if (this.isBottomIntersect) {
        this.setHasNewMessage(false)
      }
    },
    handleMutate() {
      if (this.isBottomIntersect) {
        this.scrollToBottom()
      } else {
        this.setHasNewMessage(true)
      }
    },
    scrollToBottom() {
      const bottom = this.$refs.bottom
      this.$vuetify.goTo(bottom, { duration: 0 })
    }
  }
};
</script>

