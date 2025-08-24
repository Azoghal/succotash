interface ICardProps {
    title: string;
    body: string;
}

export default function Card(props: ICardProps): JSX.Element {
    return (
        <div className="card">
            <div className="card__title">{props.title}</div>
            <div className="card__body">{props.body}</div>
            <div className="card__quick-action">
                <button className="c-btn c-btn__alternate">Quick Action</button>
            </div>
        </div>
    );
}
