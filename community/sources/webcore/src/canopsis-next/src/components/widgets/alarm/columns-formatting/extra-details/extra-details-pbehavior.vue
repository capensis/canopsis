<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          :style="{ backgroundColor: color }"
          class="c-extra-details__badge"
          v-on="on"
        >
          <v-icon
            :color="iconColor"
            small
          >
            {{ pbehaviorInfo.icon_name }}
          </v-icon>
        </span>
      </template>
      <div>
        <strong>{{ $t('alarm.actions.iconsTitles.pbehaviors') }}</strong>
        <div>
          <div class="mt-2 font-weight-bold">
            {{ pbehaviorInfo.name }}
          </div>
          <div v-if="pbehaviorInfo.author">
            {{ $t('common.author') }}: {{ pbehaviorInfo.author }}
          </div>
          <div v-if="pbehaviorInfo.type_name">
            {{ $t('common.type') }}: {{ pbehaviorInfo.type_name }}
          </div>
          <div v-if="pbehaviorInfo.reason_name">
            {{ $t('common.reason') }}: {{ pbehaviorInfo.reason_name }}
          </div>
          <div v-if="pbehaviorInfo.last_comment">
            {{ $t('alarm.fields.lastComment') }}:
            <div class="ml-2">
              -&nbsp;
              <template v-if="pbehaviorInfo.last_comment.author">
                {{ pbehaviorInfo.last_comment.author }}:&nbsp;
              </template>{{ pbehaviorInfo.last_comment.message }}
            </div>
          </div>
          <v-divider />
        </div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { getMostReadableTextColor } from '@/helpers/color';
import { getPbehaviorColor } from '@/helpers/entities/pbehavior/form';

export default {
  props: {
    pbehaviorInfo: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    color() {
      return getPbehaviorColor(this.pbehaviorInfo);
    },

    iconColor() {
      return getMostReadableTextColor(this.color, { level: 'AA', size: 'large' });
    },
  },
};
</script>
