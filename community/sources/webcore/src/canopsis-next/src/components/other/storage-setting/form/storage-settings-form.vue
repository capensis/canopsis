<template lang="pug">
  v-layout(column)
    c-information-block(
      :title="$t('storageSetting.alarm.title')",
      :help-text="$t('storageSetting.alarm.titleHelp')",
      help-icon-color="info"
    )
      template(v-if="history.alarm", #subtitle="") {{ alarmSubTitle }}
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
      :help-text="$t('storageSetting.entity.titleHelp')",
      help-icon-color="info"
    )
      template(v-if="history.entity", #subtitle="") {{ entitySubTitle }}
      v-radio-group(v-field="form.entity.archive", hide-details, mandatory, row)
        v-radio(:value="true", :label="$t('storageSetting.entity.archiveEntity')", color="primary")
        v-radio(:value="false", :label="$t('storageSetting.entity.deleteEntity')", color="primary")
      v-checkbox(
        v-field="form.entity.archive_dependencies",
        :label="$t('storageSetting.entity.archiveDependencies')",
        color="primary"
      )
        template(#append="")
          c-help-icon(:text="$t('storageSetting.entity.archiveDependenciesHelp')", color="info", max-width="300", top)
      v-flex
        v-btn.primary.ma-0.mb-4(@click="$emit('clean-entities')") {{ $t('storageSetting.entity.cleanStorage') }}
    c-information-block(
      :title="$t('storageSetting.remediation.title')",
      help-icon-color="info"
    )
      template(v-if="history.remediation", #subtitle="") {{ remediationSubTitle }}
      c-enabled-duration-field(
        v-field="form.remediation.delete_after",
        :label="$t('storageSetting.remediation.deleteAfter')",
        :help-text="$t('storageSetting.remediation.deleteAfterHelpText')",
        :name="remediationDeleteAfterFieldName"
      )
      c-enabled-duration-field(
        v-field="form.remediation.delete_stats_after",
        :label="$t('storageSetting.remediation.deleteStatsAfter')",
        :help-text="$t('storageSetting.remediation.deleteStatsAfterHelpText')",
        :name="remediationDeleteStatsAfterFieldName"
      )
      c-enabled-duration-field(
        v-field="form.remediation.delete_mod_stats_after",
        :label="$t('storageSetting.remediation.deleteModStatsAfter')",
        :help-text="$t('storageSetting.remediation.deleteModStatsAfterHelpText')",
        :name="remediationDeleteModStatsAfterFieldName"
      )
    c-information-block(
      :title="$t('storageSetting.pbehavior.title')",
      help-icon-color="info"
    )
      template(v-if="history.pbehavior", #subtitle="") {{ pbehaviorSubTitle }}
      c-enabled-duration-field(
        v-field="form.pbehavior.delete_after",
        :label="$t('storageSetting.pbehavior.deleteAfter')",
        :help-text="$t('storageSetting.pbehavior.deleteAfterHelpText')",
        :name="pbehaviorDeleteAfterFieldName"
      )
    c-information-block(
      :title="$t('storageSetting.junit.title')",
      help-icon-color="info"
    )
      template(v-if="history.junit", #subtitle="") {{ junitSubTitle }}
      c-enabled-duration-field(
        v-field="form.junit.delete_after",
        :label="$t('storageSetting.junit.deleteAfter')",
        :help-text="$t('storageSetting.junit.deleteAfterHelpText')",
        :name="junitDeleteAfterFieldName"
      )
    c-information-block(
      :title="$t('storageSetting.healthCheck.title')",
      help-icon-color="info"
    )
      template(v-if="history.health_check", #subtitle="") {{ healthCheckSubTitle }}
      c-enabled-duration-field(
        v-field="form.health_check.delete_after",
        :label="$t('storageSetting.healthCheck.deleteAfter')",
        :name="healthCheckDeleteAfterFieldName"
      )
    c-information-block(
      :title="$t('storageSetting.webhook.title')",
      :help-text="$t('storageSetting.webhook.titleHelp')",
      help-icon-color="info"
    )
      template(v-if="history.webhook", #subtitle="") {{ webhookSubTitle }}
      c-enabled-duration-field(
        v-field="form.webhook.delete_after",
        :label="$t('storageSetting.webhook.deleteAfter')",
        :help-text="$t('storageSetting.webhook.deleteAfterHelpText')",
        :name="webhookDeleteAfterFieldName"
      )
      c-enabled-field(
        v-field="form.webhook.log_credentials",
        :label="$t('storageSetting.webhook.logCredentials')",
        :name="webhookLogCredentialsFieldName",
        hide-details
      )
        template(#append="")
          c-help-icon(:text="$t('storageSetting.webhook.logCredentialsHelpText')", color="info", top)
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

    remediationDeleteAfterFieldName() {
      return 'remediation.delete_after';
    },

    remediationDeleteStatsAfterFieldName() {
      return 'remediation.delete_stats_after';
    },

    remediationDeleteModStatsAfterFieldName() {
      return 'remediation.delete_mod_stats_after';
    },

    pbehaviorDeleteAfterFieldName() {
      return 'pbehavior.delete_after';
    },

    healthCheckDeleteAfterFieldName() {
      return 'health_check.delete_after';
    },

    webhookDeleteAfterFieldName() {
      return 'webhook.delete_after';
    },

    webhookLogCredentialsFieldName() {
      return 'webhook.log_credentials';
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

    webhookSubTitle() {
      return this.$t('storageSetting.history.scriptLaunched', {
        launchedAt: convertDateToString(this.history.webhook),
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
