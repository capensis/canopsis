<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
      disable-resize
    >
      <template #activator="{ on }">
        <c-alarm-pbehavior-chip
          :color="pbehaviorInfo.color"
          :icon="pbehaviorInfo.icon_name"
          v-on="on"
        />
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
            {{ $t('common.lastComment') }}:
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
export default {
  props: {
    pbehaviorInfo: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
