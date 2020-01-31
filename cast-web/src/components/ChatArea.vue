<template>
  <v-container class="bg">
    <v-layout wrap>
      <template v-if="messages.length !== 0">
        <MessageItem
          v-for="msg in messages"
          :key="msg.id"
          :message="msg"
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
    </v-layout>
  </v-container>
</template>

<script>
import MessageItem from './MessageItem';

export default {
  name: 'ChatArea',
  components: {
    MessageItem,
  },
  props: {
    messages: {
      type: Array
    }
  },
  data: () => ({

  }),


};
</script>

