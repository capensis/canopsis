<template>
  <c-layout column>
    <c-layout gap="5">
      <c-enabled-field
        :model-value="query.group"
        :label="$t('alarm.timeline.groupItems')"
        @update:model-value="updateGroup"
      />
      <c-enabled-field
        :model-value="isCommentType"
        :label="$t('alarm.timeline.onlyComments')"
        @update:model-value="updateOnlyComments"
      />
      <v-spacer />
      <v-btn
        v-if="hasExportAccess"
        :loading="downloading"
        color="primary"
        size="small"
        variant="outlined"
        @click="exportTimeline"
      >
        <v-icon class="mr-1" size="small">
          file_download
        </v-icon>
        <span>{{ $t('alarm.timeline.export') }}</span>
      </v-btn>
    </c-layout>
    <alarm-timeline-days :days="days" :is-html-enabled="isHtmlEnabled" />
    <c-pagination
      :total="meta.total_count"
      :limit="meta.per_page"
      :page="meta.page"
      @update:page="updatePage"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { ALARM_STEPS_TYPES, USER_PERMISSIONS } from '@/constants';

import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';
import { getAlarmTimelineExportDownloadFileUrl } from '@/helpers/entities/alarm/url';
import { openUrlInNewTab } from '@/helpers/url';

import { useCanPermission } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useFilePollingWithPending } from '@/hooks/polling';
import { usePopups } from '@/hooks/popups';
import { useAlarm } from '@/hooks/store/modules/alarm';

import AlarmTimelineDays from './alarm-timeline-days.vue';

export default {
  components: { AlarmTimelineDays },

  props: {
    alarmId: {
      type: String,
      default: '',
    },
    steps: {
      type: Array,
      default: () => [],
    },
    meta: {
      type: Object,
      default: () => ({}),
    },
    query: {
      type: Object,
      default: () => ({}),
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:query'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const popups = usePopups();
    const { createAlarmTimelineExport, fetchAlarmTimelineExport } = useAlarm();
    const { hasAccess: hasExportAccess } = useCanPermission(
      USER_PERMISSIONS.business.alarmsList.actions.exportAlarmTimeline,
    );

    const { pending: downloading, poll: generateFile } = useFilePollingWithPending({
      createHandler: async () => {
        const exportData = await createAlarmTimelineExport({
          data: { alarm: props.alarmId },
        });

        return { _id: exportData._id ?? exportData.id };
      },
      fetchHandler: fetchAlarmTimelineExport,
      endHandler: fileData =>
        openUrlInNewTab(getAlarmTimelineExportDownloadFileUrl(fileData._id ?? fileData.id)),
    });

    const days = computed(() => groupAlarmSteps(props.steps));
    const isCommentType = computed(() => props.query.type === ALARM_STEPS_TYPES.comment);

    const updateGroup = group => emit('update:query', { ...props.query, group, page: 1 });
    const updateOnlyComments = (onlyComments) => {
      const { type, ...newQuery } = props.query;

      newQuery.page = 1;

      if (onlyComments) {
        newQuery.type = ALARM_STEPS_TYPES.comment;
      }

      emit('update:query', newQuery);
    };

    const updatePage = page => emit('update:query', { ...props.query, page });

    const exportTimeline = async () => {
      if (downloading.value) {
        return;
      }

      try {
        await generateFile();
      } catch {
        popups.error({ text: t('alarm.popups.timelineExportFailed') });
      }
    };

    return {
      downloading,
      hasExportAccess,
      days,
      isCommentType,
      exportTimeline,
      updateGroup,
      updateOnlyComments,
      updatePage,
    };
  },
};
</script>
