<template>
  <div class="test-suite-message">
    <v-card>
      <v-card-title
        class="white--text"
        :style="{ backgroundColor: color }"
      >
        {{ title }}
      </v-card-title>
      <v-card-text class="pt-1">
        <system-message
          :value="value"
          :label="label"
          :file-name="fileName"
        />
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import SystemMessage from './system-message.vue';

export default {
  components: { SystemMessage },
  props: {
    value: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: false,
    },
    fileName: {
      type: String,
      default: 'message',
    },
    color: {
      type: String,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.test-suite-message {
  & ::v-deep .v-sheet {
    border-radius: 5px;
    overflow: hidden;
  }
}
</style>
