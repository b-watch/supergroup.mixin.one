<template>
  <div class="card">
    <van-cell
      :title="invitation.code"
      :value="value"
      :label="label"
      :value-class="valueClass"
      @click="copy"
    />
  </div>
</template>

<script>
export default {
  name: "InvitationCodeListItem",
  props: ["invitation"],
  computed: {
    status() {
      if (this.invitation.is_used && this.invitation.invitee) {
        if (this.invitation.invitee.state == "pending") {
          return "pending";
        } else if (this.invitation.invitee.state == "paid") {
          return "used";
        }
      } else {
        return "unused";
      }
    },
    value() {
      switch (this.status) {
        case "pending":
          return this.$t("invitation.status_pending");
        case "used":
          return this.$t("invitation.status_used");
        case "unused":
          return this.$t("invitation.status_unused");
      }
    },
    valueClass() {
      switch (this.status) {
        case "pending":
          return "pending";
        case "unused":
          return "unused";
        case "unused":
          return "";
      }
    },
    label() {
      switch (this.status) {
        case "pending":
          return this.invitation.invitee.full_name + this.$t("invitation.used");
        case "used":
          return this.invitation.invitee.full_name + this.$t("invitation.used");
        case "unused":
          return this.$t("invitation.click_copy");
      }
    }
  },
  methods: {
    copy() {
      if (this.status == "unused") {
        this.$copyText(this.invitation.code);
        this.$toast(this.$t("invitation.copied"));
      }
    }
  }
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

.unused {
  color: #52c41a;
}

.pending {
  color: #f5222d;
}
</style>