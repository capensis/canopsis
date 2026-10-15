import { exportAlarmToPdf } from '@/helpers/file/pdf';

const { compileMock, htmlMock, saveMock } = vi.hoisted(() => ({
  compileMock: vi.fn().mockResolvedValue('<div>Alarm</div>'),
  htmlMock: vi.fn(),
  saveMock: vi.fn().mockResolvedValue(undefined),
}));

vi.mock('jspdf', () => ({
  default: class JsPDFMock {
    html(...args) {
      htmlMock(...args);

      return { save: saveMock };
    }
  },
}));

vi.mock('@/helpers/handlebars', () => ({
  compile: compileMock,
}));

vi.mock('@/helpers/handlebars/alarm-export-pdf-helper', () => ({
  createInstanceWithHelpers: vi.fn(() => ({})),
}));

describe('PDF file helpers', () => {
  afterEach(() => {
    document.body.innerHTML = '';
    vi.clearAllMocks();
  });

  it('excludes the application DOM while preserving the PDF overlay', async () => {
    document.body.innerHTML = `
      <main id="app"><section class="nested-content"></section></main>
      <div class="v-overlay-container"></div>
      <div class="html2pdf__overlay"></div>
    `;

    await exportAlarmToPdf('<div>Alarm</div>', { _id: 'alarm-id' });

    const [, options] = htmlMock.mock.calls[0];
    const { ignoreElements } = options.html2canvas;

    expect(ignoreElements(document.querySelector('#app'))).toBe(true);
    expect(ignoreElements(document.querySelector('.v-overlay-container'))).toBe(true);
    expect(ignoreElements(document.querySelector('.html2pdf__overlay'))).toBe(false);
    expect(ignoreElements(document.querySelector('.nested-content'))).toBe(false);
    expect(saveMock).toHaveBeenCalledWith('alarm-alarm-id.pdf');
  });
});
