const pollSvcEndpoint = 'http://localhost:17000';

async function initPoll() {
  let talkId = null;
  try {
    talkId = await (await fetch(`${pollSvcEndpoint}/config/current`)).text();
  } catch { /* ignore. */ }

  if (talkId) {
    const slideEls = document.querySelectorAll('section.slides > article');
    const submittedLabels = {};
    window.handleSlideUpdate = curSlide => {
      if (curSlide >= 0 && curSlide < slideEls.length) {
        const labels = slideEls[curSlide].getElementsByClassName('poll-label');

        if (labels.length === 1) {
          const name = labels[0].innerText;
          console.log(submittedLabels);
          if (submittedLabels[name]) return;

          fetch(`${pollSvcEndpoint}/v1/labels`, {
            method: 'POST',
            body: JSON.stringify({
              talk_name: talkId,
              name,
              timestamp: new Date().toISOString(),
            })
          }).then(() => submittedLabels[name] = true)
        }
      }
    };
  }
}

initPoll().then();
