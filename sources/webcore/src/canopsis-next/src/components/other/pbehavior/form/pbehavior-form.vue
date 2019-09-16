<template lang="pug">
  div
    pbehavior-general-form(v-model="form.general")
    pbehavior-exdates-form.mt-2(
    v-show="form.general.rrule",
    v-model="form.exdate",
    @input="updateField('exdate', $event)"
    )
    pbehavior-comments-form.mt-2(v-model="form.comments", @input="updateField('comments', $event)", :author="author")
</template>

<script>
import formMixin from '@/mixins/form';

import PbehaviorGeneralForm from './partials/pbehavior-general-form.vue';
import PbehaviorExdatesForm from './partials/pbehavior-exdates-form.vue';
import PbehaviorCommentsForm from './partials/pbehavior-comments-form.vue';

export default {
  components: {
    PbehaviorGeneralForm,
    PbehaviorExdatesForm,
    PbehaviorCommentsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    author: {
      type: String,
    },
  },
};
</script>
