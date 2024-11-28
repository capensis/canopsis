<template>
  <c-action-btn
    :disabled="pending"
    :loading="pending"
    :tooltip="$t('common.download')"
    icon="file_download"
    @click="download"
  />
</template>

<script>
import { computed } from 'vue';

import { saveFile } from '@/helpers/file/files';

import { usePendingHandler } from '@/hooks/query/pending';
import { useDbExport } from '@/hooks/store/modules/db-export';

export default {
  props: {
    id: {
      type: String,
      required: false,
    },
    ids: {
      type: Array,
      default: () => [],
    },
    pbehavior: {
      type: Boolean,
      default: false,
    },
    eventFilter: {
      type: Boolean,
      default: false,
    },
    linkRule: {
      type: Boolean,
      default: false,
    },
    idleRule: {
      type: Boolean,
      default: false,
    },
    flappingRule: {
      type: Boolean,
      default: false,
    },
    scenario: {
      type: Boolean,
      default: false,
    },
    dynamicInfo: {
      type: Boolean,
      default: false,
    },
    declareTicket: {
      type: Boolean,
      default: false,
    },
    metaAlarmRule: {
      type: Boolean,
      default: false,
    },
    instruction: {
      type: Boolean,
      default: false,
    },
    job: {
      type: Boolean,
      default: false,
    },
    jobConfig: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const {
      exportPbehaviorsDb,
      exportEventFiltersDb,
      exportLinkRulesDb,
      exportIdleRulesDb,
      exportFlappingRulesDb,
      exportScenariosDb,
      exportDynamicInfosDb,
      exportDeclareTicketsDb,
      exportMetaAlarmRulesDb,
      exportInstructionsDb,
      exportJobsDb,
      exportJobConfigsDb,
    } = useDbExport();

    const action = computed(() => ({
      [props.pbehavior]: exportPbehaviorsDb,
      [props.eventFilter]: exportEventFiltersDb,
      [props.linkRule]: exportLinkRulesDb,
      [props.idleRule]: exportIdleRulesDb,
      [props.flappingRule]: exportFlappingRulesDb,
      [props.scenario]: exportScenariosDb,
      [props.dynamicInfo]: exportDynamicInfosDb,
      [props.declareTicket]: exportDeclareTicketsDb,
      [props.metaAlarmRule]: exportMetaAlarmRulesDb,
      [props.instruction]: exportInstructionsDb,
      [props.job]: exportJobsDb,
      [props.jobConfig]: exportJobConfigsDb,
    }.true));

    const {
      pending,
      handler: download,
    } = usePendingHandler(async () => {
      const { headers = {}, data } = await action.value?.({
        ids: props.id ? [props.id] : props.ids,
      }) ?? {};

      const filename = String(headers['content-disposition']).replace(/^.*filename="|"$/g, '');

      return saveFile(data, filename);
    });

    return {
      pending,

      download,
    };
  },
};
</script>
