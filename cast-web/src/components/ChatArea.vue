<template>
  <v-container
    ref="chatArea"
  >
    <div
      v-if="messages.length !== 0"
      v-mutate.child="handleMutate"
    >
      <message-bubble
        v-for="(msg, index) in messages"
        :key="msg.id"
        :idx="index"
        :message="msg"
        :pip="true"
        :prev="getPrev(index, messages)"
      />
    </div>
    <template v-else>
      <p
        class="text-center font-weight-light mt-4"
        style="width: 100%;"
      >
        No messages
      </p>
    </template>
    <message-viewer />
    <live-floating-message />
    <div
      ref="bottom"
      v-intersect="handleBottomIntersect"
    />
  </v-container>
</template>

<script>
import { mapState } from 'vuex'
import MessageViewer from './message/MessageViewer'
import LiveFloatingMessage from './message/LiveFloatingMessage'
import MessageBubble from './message/MessageBubble'
import { mapMutations } from 'vuex'

export default {
  name: 'ChatArea',
  components: {
    MessageBubble,
    MessageViewer,
    LiveFloatingMessage
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
    handleBottomIntersect(entries, observer) {
      this.isBottomIntersect = entries[0].isIntersecting
      if (this.isBottomIntersect) {
        this.setHasNewMessage(false)
      }
    },
    handleMutate(data) {
      if (this.isBottomIntersect) {
        this.scrollToBottom()
      } else {
        console.log(data)
        this.setHasNewMessage(true)
      }
    },
    scrollToBottom() {
      const content = document.querySelector('#chatContent')
      content.scrollTop = content.scrollHeight
      this.setHasNewMessage(false)
    }
  }
};
</script>

