<template>
  <v-container>
    <template v-if="messages.length !== 0">
      <MessageItem
        v-for="(msg, index) in messages"
        :key="msg.id"
        :message="msg"
        :prev="getPrev(index, messages)"
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
  </v-container>
</template>

<script>
import MessageItem from './message/MessageItem';
import { mapState } from 'vuex'

export default {
  name: 'ChatArea',
  components: {
    MessageItem,
  },
  computed: {
    ...mapState('message', {
      messages: state => state.messages
    })
  },
  methods: {
    getPrev(index, messages) {
      return index === 0 ? null : messages[index-1]
    }
  }
};
</script>

