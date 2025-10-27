<script>
import VTextarea from 'vuetify/lib/components/VTextarea';

import { getStringLinesCount } from '@/helpers/string';

export default {
  extends: VTextarea,
  methods: {
    calculateInputHeight() {
      const { input } = this.$refs;

      if (!input) return;
      input.style.height = '0';

      const computedStyle = window.getComputedStyle(input);
      const lineHeight = parseFloat(computedStyle.lineHeight);
      const linesCount = getStringLinesCount(this.value);
      const height = Math.ceil(lineHeight * linesCount);
      const minHeight = parseInt(this.rows, 10) * parseFloat(this.rowHeight);

      input.style.height = `${Math.max(height, minHeight)}px`;
    },
  },
};
</script>
