const pollSvcEndpoint = 'http://localhost:17000';

async function initPoll() {
  let talkId = null;
  try {
    talkId = await (await fetch(`${pollSvcEndpoint}/config/current`)).text();
  } catch { /* ignore. */ }

  if (talkId) {
    const slideEls = document.querySelectorAll('section.slides > article');
    window.handleSlideUpdate = curSlide => {
      if (curSlide >= 0 && curSlide < slideEls.length) {
        const labels = slideEls[curSlide].getElementsByClassName('poll-label');

        if (labels.length === 1) {
          fetch(`${pollSvcEndpoint}/v1/labels`, {
            method: 'POST',
            body: JSON.stringify({
              talk_name: talkId,
              name: labels[0].innerText,
              timestamp: new Date().toISOString(),
            })
          })
        }
      }
    };
  }
}

initPoll().then();
