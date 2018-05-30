<template lang="pug">
  div
    span(:class="[`bg-${color}`, 'badge']") {{ text }}
    v-icon(:color="color") account_circle
</template>

<script>
export default {
  props: {
    stateId: {
      required: true,
    },
  },
  computed: {
    color() {
      switch (this.stateId) {
        case 0:
          return 'green';
        case 1:
          return 'yellow';
        case 2:
          return 'orange';
        case 3:
          return 'red';
        default:
          return 'purple';
      }
    },
    text() {
      switch (this.stateId) {
        case 0:
          return 'Info';
        case 1:
          return 'Minor';
        case 2:
          return 'Major';
        case 3:
          return 'Critical';
        default:
          return 'Unknown';
      }
    },
  },
};
</script>

<style scoped>
  .badge {
    display: inline-block;
    min-width: 10px;
    padding: 3px 7px;
    font-size: 12px;
    font-weight: 700;
    line-height: 1;
    color: #fff;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    background-color: #777;
    border-radius: 10px;
  }

  .badge.bg-green {
    background-color: #00a65a
  }

  .badge.bg-red {
    background-color: #CF0000
  }

  .badge.bg-yellow {
    background-color: #FFE000
  }

  .badge.bg-orange {
    background-color: #FF9900
  }
</style>
