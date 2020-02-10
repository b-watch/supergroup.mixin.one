<template>
  <v-container
    ref="chatArea"
    v-mutate.child="handleMutate"
  >
    <template v-if="messages.length !== 0">
      <message-bubble
        v-for="(msg, index) in messages"
        :key="msg.id"
        :idx="index"
        :message="msg"
        :pip="true"
        :prev="getPrev(index, messages)"
      />
    </template>
    <template v-else>
      <p
        class="text-center font-weight-light mt-4"
        style="width: 100%;"
      >
        还没有开始广播。
      </p>
    </template>
    <message-viewer />
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
      console.log('scroll bottom')
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
    handleMutate() {
      if (this.isBottomIntersect) {
        this.scrollToBottom()
      } else {
        this.setHasNewMessage(true)
      }
    },
    scrollToBottom() {
      const content = document.querySelector('#chatContent')
      content.scrollTop = content.scrollHeight
    }
  }
};
</script>

