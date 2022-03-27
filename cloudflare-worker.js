addEventListener("fetch", event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  const url = new URL(request.url);

  url.protocol = "https:";
  url.hostname = "us-central1-my-cloud-giulio.cloudfunctions.net";
  url.pathname = "/function-1";
  accessToken = gcloudbearer

  return fetch(new Request(url.toString(), new Request(request, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }))
  );
}


