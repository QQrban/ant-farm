window.onload = function () {
  fetch("/paths")
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((data) => {
      const vertices = data[0].Vertices;
      const edgeStrings = data[0].Edges;
      const scaleFactor = 50;
      const cytoscapeElements = vertices
        .map((vertex) => ({
          data: { id: vertex.Name },
          position: {
            x: vertex.Position.X * scaleFactor,
            y: vertex.Position.Y * scaleFactor,
          },
        }))
        .concat(
          edgeStrings.map((edgeString) => {
            const [source, target] = edgeString.split("-");
            return {
              data: {
                id: source + "-" + target,
                source: source,
                target: target,
              },
            };
          })
        );

      const cy = cytoscape({
        container: document.getElementById("cy"),
        elements: cytoscapeElements,
        style: [
          {
            selector: "node",
            css: {
              label: "data(id)",
              "text-valign": "center",
              "text-halign": "center",
              height: "45px",
              width: "45px",
              color: "#fff",
            },
          },
          {
            selector: "edge",
            css: {
              "target-arrow-shape": "triangle",
            },
          },
          {
            selector: ".background",
            css: {
              "text-background-opacity": 1,
              color: "#fff",
              "text-background-color": "#000",
            },
          },
          {
            selector: ".outline",
            css: {
              color: "#fff",
              "text-outline-color": "#000",
              "text-outline-width": 3,
            },
          },
          {
            selector: ".top-center",
            style: {
              "text-valign": "top",
              "text-halign": "center",
            },
          },
          {
            selector: 'node[id ^= "ant"]',
            css: {
              "background-image": "url('./ant.svg')",
              "background-opacity": 0,
              label: "",
              "text-valign": "center",
              "text-halign": "center",
              color: "#fff",
              width: "20px",
              height: "20px",
            },
          },
        ],
        layout: {
          name: "preset",
        },
      });
      cy.nodes().ungrabify();

      let currentPathIndex = 0;
      let currentVertexIndex = 0;
      let isAnimating = false;

      function animatePath() {
        const pathColors = [
          "blue",
          "red",
          "green",
          "purple",
          "orange",
          "dodgerblue",
          "teal",
          "pink",
          "brown",
          "gray",
          "indigo",
          "violet",
          "cyan",
          "magenta",
          "lime",
          "tan",
          "olive",
          "navy",
          "maroon",
          "coral",
          "royalblue",
          "mediumblue",
          "darkorange",
          "darkred",
          "darkgreen",
          "darkblue",
          "darkviolet",
          "darkgrey",
          "lightblue",
          "lightgreen",
          "lightcoral",
          "lightpink",
          "lightsalmon",
          "lightseagreen",
          "lightskyblue",
          "lightsteelblue",
          "lightyellow",
          "limegreen",
          "mediumaquamarine",
          "mediumorchid",
          "mediumpurple",
          "mediumseagreen",
          "mediumslateblue",
          "mediumspringgreen",
          "mediumturquoise",
          "mediumvioletred",
          "midnightblue",
          "mintcream",
          "mistyrose",
          "moccasin",
          "oldlace",
          "olivedrab",
          "orangered",
          "orchid",
          "palegoldenrod",
          "palegreen",
          "paleturquoise",
          "palevioletred",
          "papayawhip",
          "peachpuff",
          "peru",
          "plum",
          "powderblue",
        ];
        const color = pathColors[currentPathIndex % pathColors.length];

        if (!isAnimating) return;

        const path = data[0].Paths[currentPathIndex];
        const vertices = path.trim().split(" ");

        if (currentVertexIndex < vertices.length) {
          const vertex = vertices[currentVertexIndex];
          const nextVertex = vertices[currentVertexIndex + 1];

          let edgeId = vertex + "-" + nextVertex;
          let edge = cy.getElementById(edgeId);

          if (edge.length === 0) {
            edgeId = nextVertex + "-" + vertex;
            edge = cy.getElementById(edgeId);
          }

          cy.getElementById(vertex).style("background-color", color);
          cy.getElementById(edgeId).style("line-color", color);

          currentVertexIndex++;
          setTimeout(animatePath, 100);
        } else {
          currentPathIndex++;

          document.querySelector(".path-counter").innerHTML = currentPathIndex;
          currentVertexIndex = 0;
          if (currentPathIndex < data[0].Paths.length) {
            setTimeout(animatePath, 100);
          } else {
            isAnimating = false;
          }
        }
      }

      function startAnimation() {
        isAnimating = true;
        animatePath();
      }
      function stopAnimation() {
        isAnimating = false;
        currentPathIndex = 0;
        currentVertexIndex = 0;

        document.querySelector(".path-counter").innerHTML = 0;
        cy.nodes().style("background-color", "");
        cy.edges().style("line-color", "");
      }

      const startBtn = document.getElementById("start-button");
      const stopBtn = document.getElementById("stop-button");
      startBtn.addEventListener("click", () => {
        startAnimation();
      });
      stopBtn.addEventListener("click", stopAnimation);

      const antMoves = data[0].AntMoves;
      const antsNumber = Math.max(
        ...antMoves.flat(2).map((item) => parseInt(item.split("-")[0]))
      );

      const start = data[0].Start;
      const ants = [];
      const startVertex = cy.getElementById(start);

      for (let i = 1; i <= antsNumber; i++) {
        const ant = {
          group: "nodes",
          data: { id: "ant" + i },
          position: {
            x: startVertex.position("x"),
            y: startVertex.position("y"),
          },
          classes: "ant",
        };

        ants.push(ant);
      }

      cy.add(ants);
      cy.nodes('[id ^= "ant"]').ungrabify();

      let currentMoveIndex = 0;
      let animationTimeoutId;
      let isForbidden = false;

      function animateAntsMoves(moveIndex) {
        const currentMoves = antMoves[moveIndex];
        currentMoves.forEach((move) => {
          const [antId, targetVertexId] = move.split("-");
          const ant = cy.getElementById("ant" + antId);
          const targetVertex = cy.getElementById(targetVertexId);

          ant.animate({
            position: targetVertex.position(),
            duration: 1000,
          });
        });
      }

      function startAntsAnimation() {
        if (currentMoveIndex < antMoves.length) {
          animateAntsMoves(currentMoveIndex);
          currentMoveIndex++;
          animationTimeoutId = setTimeout(startAntsAnimation, 1100);
        }
      }

      function resetAntsPosition() {
        const startVertex = cy.getElementById(start);
        const startPosition = startVertex.position();

        cy.nodes('[id ^= "ant"]').stop().position(startPosition);
        currentMoveIndex = 0;
        isForbidden = false;
        clearTimeout(animationTimeoutId);
      }

      const resetAntsBtn = document.getElementById("reset-ants");
      resetAntsBtn.addEventListener("click", resetAntsPosition);

      const sendAntsBtn = document.getElementById("send-ants");
      sendAntsBtn.addEventListener("click", () => {
        if (!isForbidden) {
          startAntsAnimation();
          isForbidden = true;
        }
      });

      document.querySelectorAll('input[name="antImg"]').forEach((radio) => {
        radio.addEventListener("change", function () {
          const selectedAntImage = this.value;
          cy.style()
            .selector('node[id ^= "ant"]')
            .css({
              "background-image": `url('./${selectedAntImage}')`,
              // ... оставшиеся стили
            })
            .update(); // обновить стили в графе
        });
      });
    })
    .catch((error) => {
      console.error("Error:", error);
    });
};
