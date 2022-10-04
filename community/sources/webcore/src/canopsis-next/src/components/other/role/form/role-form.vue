<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name")
    v-text-field(
      v-field="form.description",
      :label="$t('common.description')",
      data-test="description"
    )

    c-information-block(:title="$t('role.expirationSettings')")
      c-enabled-field(v-model="form.auth_config.intervals_enabled")

      v-layout(v-if="form.auth_config.intervals_enabled", row)
        c-information-block(
          :title="$t('role.inactivityInterval')",
          :help-text="$t('role.inactivityIntervalHelpText')"
        )
          c-duration-field(v-field="form.auth_config.inactivity_interval", long)
        c-information-block.ml-3(
          :title="$t('role.expirationInterval')",
          :help-text="$t('role.expirationIntervalHelpText')"
        )
          c-duration-field(v-field="form.auth_config.expiration_interval", long)
    view-selector(v-field="form.defaultview")
</template>

<script>
import ViewSelector from '@/components/forms/fields/view-selector.vue';

export default {
  inject: ['$validator'],
  components: { ViewSelector },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
