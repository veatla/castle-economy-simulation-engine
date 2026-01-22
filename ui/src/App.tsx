import { useEffect, useRef, useState } from "react";
import * as PIXI from "pixi.js";
import "./App.css";

type AgentUpdate = { id: number; x: number; z: number; type: string };

function App() {
  const [tick, setTick] = useState<number | null>(null);
  const stageRef = useRef<HTMLDivElement | null>(null);
  const appRef = useRef<PIXI.Application | null>(null);
  const spritesRef = useRef<Map<number, PIXI.Graphics>>(new Map());

  useEffect(() => {
    if (!stageRef.current) return;
    const app = new PIXI.Application({
      backgroundColor: 0x1e1e1e,
      resizeTo: stageRef.current,
      antialias: true,
      resolution: window.devicePixelRatio || 1,
    });
    stageRef.current.appendChild(app.view as HTMLCanvasElement);
    appRef.current = app;

    return () => {
      app.destroy(true, { children: true });
      appRef.current = null;
    };
  }, []);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.addEventListener("message", (ev) => {
      try {
        const data = JSON.parse(ev.data) as { tick: number; updated: AgentUpdate[] };
        setTick(data.tick);
        const app = appRef.current;
        if (!app) return;

        const W = app.renderer.width;
        const H = app.renderer.height;

        data.updated.forEach((u) => {
          if (u.type !== "agent") return;
          let g = spritesRef.current.get(u.id);
          const screenX = (u.x / 100) * W;
          const screenY = (u.z / 100) * H;

          if (!g) {
            g = new PIXI.Graphics();
            g.beginFill(0xffcc00);
            g.drawCircle(0, 0, 6);
            g.endFill();
            g.x = screenX;
            g.y = screenY;
            g.zIndex = 1;
            app.stage.addChild(g);
            spritesRef.current.set(u.id, g);
          } else {
            // smooth move
            g.x += (screenX - g.x) * 0.6;
            g.y += (screenY - g.y) * 0.6;
          }
        });

        console.log(`Agents list ${spritesRef.current.size}`);
      } catch (e) {
        // ignore
      }
    });
    return () => ws.close();
  }, []);

  return (
    <div style={{ display: "flex", gap: 12 }}>
      <div style={{ width: 500, height: 500 }} ref={stageRef} />
      <div style={{ width: 220 }}>
        <div style={{ marginBottom: 8 }}>
          <strong>Tick:</strong> {tick ?? "-"}
        </div>
        <div>
          <em>Agents rendered with Pixi.js</em>
        </div>
      </div>
    </div>
  );
}

export default App;
