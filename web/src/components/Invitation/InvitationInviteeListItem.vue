<template>
  <div class="card">
    <div class="user">
      <img class="avatar" :src="invitation.invitee.avatar_url" />
      <div class="info">
        <div class="name">{{invitation.invitee.full_name}}</div>
        <div class="mixinId">{{invitation.invitee.identity_number}}</div>
      </div>
    </div>
    <div class="code">{{this.$t("invitation.code")}}： {{invitation.code}}</div>
    <div class="time">{{this.$t("invitation.used_at")}}：{{dayjs(invitation.used_at).format('YYYY/MM/DD hh:mm')}}</div>
  </div>
</template>

<script>
import dayjs from 'dayjs'

export default {
  name: "InvitationInviteeListItem",
  props: ["invitation"],
};
</script>

<style lang="scss" scoped>
.van-cell {
  padding: 0rem;
}

.card {
  background-color: #ffffff;
  border-radius: 4px;
  box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.1);
}

.user {
  height: 2rem;
}

.info {
  float: left;
  padding-left: 1rem;
}

.name {
  font-size: 1rem;
  color: #000000;
  opacity: 0.8;
}

.mixinId {
  font-size: 0.75rem;
  color: rgba(0, 0, 0, 0.25);
}

.avatar {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  float: left;
}

.other {
  padding-top: 1rem;
}

.code {
  padding-top: 0.5rem;
  color: rgba(0, 0, 0, 0.25);
  font-size: 0.8rem;
}

.time {
  padding-top: 0.5rem;
  color: rgba(0, 0, 0, 0.25);
  font-size: 0.8rem;
}
</style>