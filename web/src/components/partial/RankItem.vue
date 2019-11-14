<template>
  <div class="rank-member-item"
    :class="itemClass"
    @click="handleMemberClick(member)">
    <div class="cell rank-member-icon">
      <a class="avatar" :href="'mixin://users/' + member.user_id">
        <img class="rank-member-icon-img" :src="member.avatar_url" />
        <img class="rank-member-icon-badge" :src="itemRankBadge" />
      </a>
    </div>
    <div class="cell rank-member-list-info">
      <div class="rank-member-name">{{ member.full_name }}</div>
    </div>
    <div class="cell rank-member-list-role">
      <div class="rank-member-usd">${{ member.tip_usd}}</div>
      <div class="rank-member-count">{{ $('reward_rank.item_count_text', {count: member.tip_count}) }}</div>
    </div>
  </div>
</template>

<script>
export default {
  props: ["member", "rank"],
  computed: {
    itemClass () {
      if (this.rank === 0) return 'first';
      else if (this.rank === 1) return 'second';
      else if (this.rank === 2) return 'third';
      else return ''
    },
    itemRankBadge () {
      if (this.rank === 0) return require('@/assets/images/rank-1st.png');
      else if (this.rank === 1) return require('@/assets/images/rank-2nd.png');
      else if (this.rank === 2) return require('@/assets/images/rank-3rd.png');
      else return ''
    }
  },
  methods: {
    handleMemberClick(member) {
      this.$emit("rank-member-click", member);
    },
  }
};
</script>

<style lang="scss" scoped>
.rank-member-item {
  display: flex;
  font-size: 12px;
  padding: 10px 15px;
  border-bottom: 1px solid #f8f8f8;
  background: #fff;
  .avatar {
    position: relative;
  }
  .rank-member-icon-badge {
    position: absolute;
    transform: rotate(45deg);
    height: 18px;
    width: auto;
    right: -8px;
    top: -6px;
  }
  .cell {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: center;
  }

  .rank-member-icon {
    padding-right: 10px;

    .rank-member-icon-img {
      width: 32px;
      height: 32px;
      border-radius: 14px;
    }
  }

  .rank-member-list-info {
    flex: 1;
  }
  .rank-member-list-role {
    align-items: flex-end;
  }

  .rank-member-name {
    height: 20px;
    font-size: 16px;
  }

  .rank-member-usd {
    font-size: 14px;
  }
}
.rank-member-item.first {
  color: #EBB62F;
}
.rank-member-item.second {
  color: #CDCDCD;
}
.rank-member-item.third {
  color: #9C4B00;
}
</style>

