<template lang="pug">
  v-layout(column)
    c-information-block(
      :title="$t('storageSettings.alarm.title')",
      :help-text="$t('storageSettings.alarm.titleHelp')"
    )
      template(v-if="history.alarm", #subtitle="") {{ alarmSubTitle }}
      c-enabled-duration-field(
        v-field="form.alarm.archive_after",
        :label="$t('storageSettings.alarm.archiveAfter')",
        :name="alarmArchiveAfterFieldName"
      )
      c-enabled-duration-field(
        v-field="form.alarm.delete_after",
        :label="$t('storageSettings.alarm.deleteAfter')",
        :name="alarmDeleteAfterFieldName"
      )
    c-information-block(
      :title="$t('storageSettings.entity.title')",
      :help-text="$t('storageSettings.entity.titleHelp')"
    )
      template(v-if="history.entity", #subtitle="") {{ entitySubTitle }}
      v-radio-group(v-field="form.entity.archive", hide-details, mandatory, row)
        v-radio(:value="true", :label="$t('storageSettings.entity.archiveEntity')", color="primary")
        v-radio(:value="false", :label="$t('storageSettings.entity.deleteEntity')", color="primary")
      v-checkbox(
        v-field="form.entity.archive_dependencies",
        :label="$t('storageSettings.entity.archiveDependencies')",
        color="primary"
      )
        template(#append="")
          c-help-icon(:text="$t('storageSettings.entity.archiveDependenciesHelp')", max-width="300", top)
      v-flex
        v-btn.primary.ma-0.mb-4(@click="$emit('clean-entities')") {{ $t('storageSettings.entity.cleanStorage') }}
    c-information-block(:title="$t('storageSettings.remediation.title')")
      template(v-if="history.remediation", #subtitle="") {{ remediationSubTitle }}
      c-enabled-duration-field(
        v-field="form.remediation.accumulate_after",
        :label="$t('storageSettings.remediation.accumulateAfter')",
        :name="remediationAccumulateAfterFieldName"
      )
      c-enabled-duration-field(
        v-field="form.remediation.delete_after",
        :label="$t('storageSettings.remediation.deleteAfter')",
        :help-text="$t('storageSettings.remediation.deleteAfterHelpText')",
        :name="remediationDeleteAfterFieldName"
      )
    c-information-block(:title="$t('storageSettings.pbehavior.title')")
      template(v-if="history.pbehavior", #subtitle="") {{ pbehaviorSubTitle }}
      c-enabled-duration-field(
        v-field="form.pbehavior.delete_after",
        :label="$t('storageSettings.pbehavior.deleteAfter')",
        :help-text="$t('storageSettings.pbehavior.deleteAfterHelpText')",
        :name="pbehaviorDeleteAfterFieldName"
      )
    c-information-block(:title="$t('storageSettings.junit.title')")
      template(v-if="history.junit", #subtitle="") {{ junitSubTitle }}
      c-enabled-duration-field(
        v-field="form.junit.delete_after",
        :label="$t('storageSettings.junit.deleteAfter')",
        :help-text="$t('storageSettings.junit.deleteAfterHelpText')",
        :name="junitDeleteAfterFieldName"
      )
    c-information-block(:title="$t('storageSettings.healthCheck.title')")
      template(v-if="history.health_check", #subtitle="") {{ healthCheckSubTitle }}
      c-enabled-duration-field(
        v-field="form.health_check.delete_after",
        :label="$t('storageSettings.healthCheck.deleteAfter')",
        :name="healthCheckDeleteAfterFieldName"
      )
</template>

<script>
import { isNumber } from 'lodash';

import { convertDateToString } from '@/helpers/date/date';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    history: {
      type: Object,
      required: true,
    },
  },
  computed: {
    junitDeleteAfterFieldName() {
      return 'junit.delete_after';
    },

    alarmArchiveAfterFieldName() {
      return 'alarm.archive_after';
    },

    alarmDeleteAfterFieldName() {
      return 'alarm.delete_after';
    },

    remediationAccumulateAfterFieldName() {
      return 'remediation.accumulate_after';
    },

    remediationDeleteAfterFieldName() {
      return 'remediation.delete_after';
    },

    pbehaviorDeleteAfterFieldName() {
      return 'pbehavior.delete_after';
    },

    healthCheckDeleteAfterFieldName() {
      return 'health_check.delete_after';
    },

    junitSubTitle() {
      return this.$t('storageSettings.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.junit),
      });
    },

    remediationSubTitle() {
      return this.$t('storageSettings.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.remediation),
      });
    },

    pbehaviorSubTitle() {
      return this.$t('storageSettings.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.pbehavior),
      });
    },

    healthCheckSubTitle() {
      return this.$t('storageSettings.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.health_check),
      });
    },

    alarmSubTitle() {
      const { time, deleted, archived } = this.history.alarm || {};

      const result = [
        this.$t('storageSettings.history.scriptLaunched', {
          launchedAt: convertDateToString(time),
        }),
      ];

      if (isNumber(deleted)) {
        result.push(this.$t('storageSettings.history.alarm.deletedCount', {
          count: deleted,
        }));
      }

      if (isNumber(archived)) {
        result.push(this.$t('storageSettings.history.alarm.archivedCount', {
          count: archived,
        }));
      }

      return result.join(' ');
    },

    entitySubTitle() {
      const { time, deleted, archived } = this.history.entity || {};

      const result = [
        this.$t('storageSettings.history.scriptLaunched', {
          launchedAt: convertDateToString(time),
        }),
      ];

      if (isNumber(deleted)) {
        result.push(this.$t('storageSettings.history.entity.deletedCount', {
          count: deleted,
        }));
      }

      if (isNumber(archived)) {
        result.push(this.$t('storageSettings.history.entity.archivedCount', {
          count: archived,
        }));
      }

      return result.join(' ');
    },
  },
  watch: {
    'form.remediation': function remediationWatcher() {
      this.$validator.validateAll([
        this.remediationAccumulateAfterFieldName,
        this.remediationDeleteAfterFieldName,
      ]);
    },
    'form.alarm': function alarmWatcher() {
      this.$validator.validateAll([
        this.alarmArchiveAfterFieldName,
        this.alarmDeleteAfterFieldName,
      ]);
    },
  },
};
</script>
