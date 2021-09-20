<template lang="pug">
  v-layout(column)
    c-information-block(
      :title="$t('storageSetting.alarm.title')",
      :help-text="$t('storageSetting.alarm.titleHelp')"
    )
      template(v-if="history.alarm", slot="subtitle") {{ alarmSubTitle }}
      c-enabled-duration-field(
        v-field="form.alarm.archive_after",
        :label="$t('storageSetting.alarm.archiveAfter')",
        :name="alarmArchiveAfterFieldName"
      )
      c-enabled-duration-field(
        v-field="form.alarm.delete_after",
        :label="$t('storageSetting.alarm.deleteAfter')",
        :name="alarmDeleteAfterFieldName"
      )
    c-information-block(
      :title="$t('storageSetting.entity.title')",
      :help-text="$t('storageSetting.entity.titleHelp')"
    )
      template(v-if="history.entity", slot="subtitle") {{ entitySubTitle }}
      v-radio-group(v-field="form.entity.archive", hide-details, mandatory, row)
        v-radio(:value="true", :label="$t('storageSetting.entity.archiveEntity')", color="primary")
        v-radio(:value="false", :label="$t('storageSetting.entity.deleteEntity')", color="primary")
      v-checkbox(
        v-field="form.entity.archive_dependencies",
        :label="$t('storageSetting.entity.archiveDependencies')",
        color="primary"
      )
        c-help-icon(
          slot="append",
          :text="$t('storageSetting.entity.archiveDependenciesHelp')",
          max-width="300",
          top
        )
      v-flex
        v-btn.primary.ma-0.mb-4(@click="$emit('clean-entities')") {{ $t('storageSetting.entity.cleanStorage') }}
    c-information-block(:title="$t('storageSetting.remediation.title')")
      template(v-if="history.remediation", slot="subtitle") {{ remediationSubTitle }}
      c-enabled-duration-field(
        v-field="form.remediation.accumulate_after",
        :label="$t('storageSetting.remediation.accumulateAfter')",
        :name="remediationAccumulateAfterFieldName"
      )
      c-enabled-duration-field(
        v-field="form.remediation.delete_after",
        :label="$t('storageSetting.remediation.deleteAfter')",
        :help-text="$t('storageSetting.remediation.deleteAfterHelpText')",
        :name="remediationDeleteAfterFieldName"
      )
    c-information-block(:title="$t('storageSetting.pbehavior.title')")
      template(v-if="history.pbehavior", slot="subtitle") {{ pbehaviorSubTitle }}
      c-enabled-duration-field(
        v-field="form.pbehavior.delete_after",
        :label="$t('storageSetting.pbehavior.deleteAfter')",
        :help-text="$t('storageSetting.pbehavior.deleteAfterHelpText')",
        :name="pbehaviorDeleteAfterFieldName"
      )
    c-information-block(:title="$t('storageSetting.junit.title')")
      template(v-if="history.junit", slot="subtitle") {{ junitSubTitle }}
      c-enabled-duration-field(
        v-field="form.junit.delete_after",
        :label="$t('storageSetting.junit.deleteAfter')",
        :help-text="$t('storageSetting.junit.deleteAfterHelpText')",
        :name="junitDeleteAfterFieldName"
      )
    c-information-block(:title="$t('storageSetting.healthCheck.title')")
      template(v-if="history.health_check", slot="subtitle") {{ healthCheckSubTitle }}
      c-enabled-duration-field(
        v-field="form.health_check.delete_after",
        :label="$t('storageSetting.healthCheck.deleteAfter')",
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
      return this.$t('storageSetting.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.junit),
      });
    },

    remediationSubTitle() {
      return this.$t('storageSetting.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.remediation),
      });
    },

    pbehaviorSubTitle() {
      return this.$t('storageSetting.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.pbehavior),
      });
    },

    healthCheckSubTitle() {
      return this.$t('storageSetting.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.health_check),
      });
    },

    alarmSubTitle() {
      const { time, deleted, archived } = this.history.alarm || {};

      const result = [
        this.$t('storageSetting.history.scriptLaunched', {
          launchedAt: convertDateToString(time),
        }),
      ];

      if (isNumber(deleted)) {
        result.push(this.$t('storageSetting.history.alarm.deletedCount', {
          count: deleted,
        }));
      }

      if (isNumber(archived)) {
        result.push(this.$t('storageSetting.history.alarm.archivedCount', {
          count: archived,
        }));
      }

      return result.join(' ');
    },

    entitySubTitle() {
      const { time, deleted, archived } = this.history.entity || {};

      const result = [
        this.$t('storageSetting.history.scriptLaunched', {
          launchedAt: convertDateToString(time),
        }),
      ];

      if (isNumber(deleted)) {
        result.push(this.$t('storageSetting.history.entity.deletedCount', {
          count: deleted,
        }));
      }

      if (isNumber(archived)) {
        result.push(this.$t('storageSetting.history.entity.archivedCount', {
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
