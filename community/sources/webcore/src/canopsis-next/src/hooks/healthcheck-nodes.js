import { useI18n } from '@/hooks/i18n';

export const useHealthcheckNodes = () => {
  const { t, te } = useI18n();

  const getNodeName = (nodeName) => {
    const nodeLabelKey = `healthcheck.nodes.${nodeName}.name`;

    return te(nodeLabelKey) ? t(nodeLabelKey) : nodeName;
  };

  const getTooltipText = (data = {}) => {
    if (data.is_unknown) {
      return t('healthcheck.unknown');
    }

    return (data.issues ?? [])
      .map((issue) => {
        const issueKey = `healthcheck.issues.${issue}`;

        return te(issueKey) ? t(issueKey, { name: getNodeName(data.name) }) : issue;
      })
      .join('\n');
  };

  return {
    getNodeName,
    getTooltipText,
  };
};
