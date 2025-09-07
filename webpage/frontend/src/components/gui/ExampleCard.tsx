import { useCallback, useState } from "react";

interface IExampleCardProps {
    title: string;
    bobOrBill: string; // TODO sort out
}

export default function ExampleCard(props: IExampleCardProps): JSX.Element {
    const [popularity, setPopularity] = useState<number>(-1);
    const [artist, setArtist] = useState("");

    const fetchData = useCallback(async () => {
    
        try {
            const response = await fetch(`http://localhost:6789/api/v1/authed/test/${props.bobOrBill}`, {credentials: "include"});
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const fetchedData = await response.json();
            console.log(props.bobOrBill,fetchedData)
        } catch (e) {
            console.error("Failed to fetch data:", e);
        }
    },[props]);   

    const getArtistPopularity = useCallback(() => {
        fetchData();

        setPopularity(20);
    }, [fetchData]);

    return (
        <div className="card">
            <div className="card__title">{props.title}</div>
            <div className="card__body">
                <div>
                    <input
                        type="text"
                        onChange={(e) => setArtist(e.target.value)}
                        value={artist}
                    />
                    {artist}
                </div>
                <div>Popularity: {popularity == -1 ? "..." : popularity}</div>
            </div>
            <div className="card__quick-action">
                <button
                    className="c-btn c-btn__alternate"
                    onClick={getArtistPopularity}
                >
                    Get Artist Popularity
                </button>
            </div>
        </div>
    );
}
