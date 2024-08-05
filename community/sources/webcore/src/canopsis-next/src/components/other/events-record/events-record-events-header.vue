<template>
  <v-layout class="gap-2" column>
    <v-layout class="gap-3 my-4" justify-center align-center>
      <span class="text-subtitle-2">{{ $t('modals.eventsRecord.subtitle', { count: eventsRecord.c }) }}</span>
      <c-action-btn
        :tooltip="$t('modals.eventsRecord.buttonTooltip')"
        type="delete"
        @click="remove"
      />
    </v-layout>
    <v-layout class="gap-3 my-4">
      <v-btn
        color="primary"
        @click="applyEventFilter"
      >
        {{ $t('eventsRecord.applyEventFilter') }}
      </v-btn>
      <v-btn
        :loading="downloading"
        color="primary"
        @click="exportJson"
      >
        <v-icon class="mr-2" color="white">
          file_download
        </v-icon>
        <span>{{ $t('common.exportToJson') }}</span>
      </v-btn>
    </v-layout>
  </v-layout>
</template>
<script>
export default {
  props: {
    eventsRecord: {
      type: Object,
      required: true,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const applyEventFilter = () => {};
    const remove = () => {};
    const exportJson = () => {};

    return {
      applyEventFilter,
      remove,
      exportJson,
    };
  },
};
</script>
