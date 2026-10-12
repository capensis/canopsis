import { mapActions } from 'pinia';

import { useEntityStore } from '@/stores/modules/entity';

export const entityPauseMixin = {
  methods: {
    ...mapActions(useEntityStore, {
      pauseEntity: 'pause',
      resumeEntity: 'resume',
      bulkPauseEntities: 'bulkPause',
      bulkResumeEntities: 'bulkResume',
    }),
  },
};
