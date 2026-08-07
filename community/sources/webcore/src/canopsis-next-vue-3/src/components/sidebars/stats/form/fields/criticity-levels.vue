<template>
  <widget-settings-item :title="$t('settings.criticityLevels')">
    <c-layout gap="2" column>
      <c-layout
        v-for="level in levels"
        :key="level"
        align="center"
        justify="space-around"
      >
        <div>{{ $t(`settings.colorsSelector.statsCriticity.${level}`) }} :</div>

        <div class="criticity-levels__field">
          <c-number-field
            :value="modelValue[level]"
            :min="mins[level]"
            :name="level"
            required
            @update:value="updateLevel(level, $event)"
          />
        </div>
      </c-layout>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { STATS_CRITICITY } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Object,
      default: () => ({
        minor: 20,
        major: 30,
        critical: 40,
      }),
    },
  },

  emits: ['update:modelValue'],

  computed: {
    levels() {
      return [
        STATS_CRITICITY.minor,
        STATS_CRITICITY.major,
        STATS_CRITICITY.critical,
      ];
    },

    mins() {
      return {
        [STATS_CRITICITY.minor]: 0,
        [STATS_CRITICITY.major]: (this.modelValue.minor ?? 0) + 1,
        [STATS_CRITICITY.critical]: (this.modelValue.major ?? 0) + 1,
      };
    },
  },

  methods: {
    updateLevel(level, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [level]: value,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.criticity-levels__field {
  width: 25%;
  min-width: 84px;
}
</style>
