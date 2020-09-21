<template lang="pug">
  v-list-group(data-test="criticityLevels")
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ $t('settings.criticityLevels') }}
    v-container
      v-layout(wrap)
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.minor') }} :
            v-flex(xs3)
              v-text-field(
                v-field.number="levels.minor",
                v-validate="'required|min_value:0'",
                :error-messages="errors.collect('minor')",
                data-test="criticityLevelsMinor",
                data-vv-name="minor",
                type="number"
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.major') }} :
            v-flex(xs3)
              v-text-field(
                v-field.number="levels.major",
                v-validate="`required|min_value:${levels.minor + 1}`",
                :error-messages="errors.collect('major')",
                data-test="criticityLevelsMajor",
                data-vv-name="major",
                type="number"
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.critical') }} :
            v-flex(xs3)
              v-text-field(
                v-field.number="levels.critical",
                v-validate="`required|min_value:${levels.major + 1}`",
                :error-messages="errors.collect('critical')",
                data-test="criticityLevelsCritical",
                data-vv-name="critical",
                type="number"
              )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'levels',
    event: 'input',
  },
  props: {
    levels: {
      type: Object,
      default: () => ({
        minor: 20,
        major: 30,
        critical: 40,
      }),
    },
  },
};
</script>
