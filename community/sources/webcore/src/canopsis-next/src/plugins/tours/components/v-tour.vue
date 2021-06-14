<script>
import VTour from 'vue-tour/src/components/VTour.vue';

export default {
  name: VTour.name,
  extends: VTour,
  methods: {
    start(startStep) {
      // Wait for the DOM to be loaded, then start the tour
      setTimeout(async () => {
        await this.customCallbacks.onStart();

        this.currentStep = typeof startStep !== 'undefined' ? parseInt(startStep, 10) : 0;
      }, this.customOptions.startTimeout);
    },

    async previousStep() {
      if (this.currentStep > 0) {
        await this.customCallbacks.onPreviousStep(this.currentStep);

        this.currentStep -= 1;
      }
    },

    async nextStep() {
      if (this.currentStep < this.numberOfSteps - 1 && this.currentStep !== -1) {
        await this.customCallbacks.onNextStep(this.currentStep);

        this.currentStep += 1;
      }
    },

    async stop() {
      await this.customCallbacks.onStop();

      document.body.classList.remove('v-tour--active');

      this.currentStep = -1;
    },

    async skip() {
      await this.customCallbacks.onSkip();

      this.stop();
    },

    async finish() {
      await this.customCallbacks.onFinish();

      this.stop();
    },
  },
};
</script>
