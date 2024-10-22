import './About.css';

function About() {
  return (
    <div className="about-container">
      <p>
        <h3>Humor Arena</h3>
        <p>It is common ground that modern LLMs are bad at humor generation. Even top-shelf models tend to memorize and repeat a few simple jokes without any originality.</p>
        <p>In our recent paper {' '}
        <a href="https://computationalcreativity.net/iccc24/papers/ICCC24_paper_128.pdf" target="_blank" rel="noopener noreferrer">
          "Humor Mechanics: Advancing Humor Generation with Multistep Reasoning"
        </a>{' '} (presented at the International Conference on Computational Creativity 2024), we show that the approach based on multistep reasoning can replicate the creativity process good enough to generate jokes which are on par with human-written jokes (with a top quality subset of "reddit jokes" dataset) according to the blind human labeling results. 
          For more details, you can read the {' '}
          <a href="https://arxiv.org/abs/2405.07280" target="_blank" rel="noopener noreferrer">full paper on arXiv</a>
          {' '}. 
          We also shared our {' '}
          <a href="https://github.com/altsoph/humor-mechanics" target="_blank" rel="noopener noreferrer">results and data</a>
          {' '} to facilitate future research.</p>
        
        <b>Humor Arena</b> is based on the research paper{' '}
        <a href="https://arxiv.org/abs/2405.07280" target="_blank" rel="noopener noreferrer">
          "Humor Mechanics: Advancing Humor Generation with Multistep Reasoning"
        </a>{' '}
        by Alexey Tikhonov and Pavel Shtykovskiy.
      </p>
      <p><strong>Now we want to go further:</strong> is there a way to improve reasoning schema? are some models more potent in terms of humor generation than others? To investigate it, we made this Humor Arena to ask people to {' '}
          <a href="https://humor.ph34r.me/arena" target="_blank" rel="noopener noreferrer">help us with blind side-by-side labeling</a>
          {' '}.
      </p>
      <p>
      <h3>Reference</h3>
      <pre className="latex-reference">
{`@article{tikhonov2024humor,
  title={Humor Mechanics: Advancing Humor Generation with Multistep Reasoning},
  author={Tikhonov, Alexey and Shtykovskiy, Pavel},
  journal={arXiv preprint arXiv:2405.07280},
  year={2024}
}`}
      </pre>
      </p>
    </div>
  );
}

export default About;
