<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.criticityLevels')}}
    v-container
      v-layout(wrap)
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.minor') }} :
            v-flex(xs3)
              v-text-field(
              type="number",
              :value="levels.minor",
              data-vv-name="minor",
              v-validate="'required|min_value:0'",
              :error-messages="errors.collect('minor')",
              @input="updateField('minor', parseInt($event, 10))",
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.major') }} :
            v-flex(xs3)
              v-text-field(
              type="number",
              :value="levels.major",
              data-vv-name="major",
              v-validate="`required|min_value:${levels.minor + 1}`",
              :error-messages="errors.collect('major')",
              @input="updateField('major', parseInt($event, 10))"
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.critical') }} :
            v-flex(xs3)
              v-text-field(
              type="number",
              :value="levels.critical",
              data-vv-name="critical",
              v-validate="`required|min_value:${levels.major + 1}`",
              :error-messages="errors.collect('critical')",
              @input="updateField('critical', parseInt($event, 10))"
              )
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
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
