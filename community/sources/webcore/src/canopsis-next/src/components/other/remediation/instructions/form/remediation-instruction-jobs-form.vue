<template lang="pug">
  v-layout.mt-2(column)
    v-layout(row)
      v-flex(v-if="!jobs.length", xs12)
        v-alert(:value="true", type="info") {{ $t('remediationInstructions.emptyJobs') }}
    draggable(
      v-field="jobs",
      :options="draggableOptions",
      :class="{ 'grey lighten-2': isDragging }",
      @start="startDragging",
      @end="endDragging"
    )
      remediation-instruction-job-field.py-1(
        v-for="(job, index) in jobs",
        v-field="jobs[index]",
        :key="job.key",
        :index="index",
        :job-number="index + 1",
        @remove="removeItemFromArray(index)"
      )
    v-layout(row, align-center)
      v-btn.ml-0(
        outline,
        :color="hasJobsErrors ? 'error' : 'primary'",
        @click="addJob"
      ) {{ $t('remediationInstructions.addJob') }}
      span.error--text(v-show="hasJobsErrors") {{ $t('remediationInstructions.errors.jobRequired') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import formArrayMixin from '@/mixins/form/array';

import RemediationInstructionJobField from './fields/remediation-instruction-job-field.vue';

export default {
  inject: ['$validator'],
  components: {
    Draggable,
    RemediationInstructionJobField,
  },
  mixins: [formArrayMixin],
  model: {
    prop: 'jobs',
    event: 'input',
  },
  props: {
    jobs: {
      type: Array,
      default: () => ([]),
    },
    name: {
      type: String,
      default: 'jobs',
    },
  },
  data() {
    return {
      isDragging: false,
    };
  },
  computed: {
    hasJobsErrors() {
      return this.errors.has(this.name);
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.job-drag-handler',
        ghostClass: 'white',
        group: {
          name: 'remediation-instruction-jobs',
          pull: false,
          put: false,
        },
      };
    },
  },
  watch: {
    jobs() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: 'min_value:1',
      getter: () => this.jobs.length,
      context: () => this,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    addJob() {
      this.addItemIntoArray({});
    },

    startDragging() {
      this.isDragging = true;
    },

    endDragging() {
      this.isDragging = false;
    },
  },
};
</script>
