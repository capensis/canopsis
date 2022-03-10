<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ $t('settings.criticityLevels') }}
    v-container
      v-layout(wrap)
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.minor') }} :
            v-flex(xs3)
              c-number-field(
                v-field="levels.minor",
                :min="0",
                name="minor",
                required
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.major') }} :
            v-flex(xs3)
              c-number-field(
                v-field="levels.major",
                :min="levels.minor + 1",
                name="major",
                required
              )
        v-flex(xs12)
          v-layout(align-center, justify-space-around)
            div {{ $t('settings.colorsSelector.statsCriticity.critical') }} :
            v-flex(xs3)
              c-number-field(
                v-field="levels.critical",
                :min="levels.major + 1",
                name="critical",
                required
              )
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

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
